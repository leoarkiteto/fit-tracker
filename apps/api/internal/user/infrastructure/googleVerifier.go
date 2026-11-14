package infrastructure

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"fittracker-api/internal/user"
)

const (
	googleDefaultJWKSURL    = "https://www.googleapis.com/oauth2/v3/certs"
	googleDefaultKeyTTL     = time.Hour
	googleMinRefreshGap     = time.Minute
	providerMaxResponseSize = 1 << 20
)

var googleMaxAgeRe = regexp.MustCompile(`(?i)max-age=(\d+)`)

// GoogleVerifier validates Google ID tokens against Google's published JWKS.
type GoogleVerifier struct {
	clientIDs []string
	jwksURL   string
	client    *http.Client
	now       func() time.Time

	mu        sync.Mutex
	keys      map[string]*rsa.PublicKey
	expiresAt time.Time
	lastFetch time.Time
	fetched   bool
}

var _ user.GoogleVerifier = (*GoogleVerifier)(nil)

// GoogleVerifierOption customises a GoogleVerifier.
type GoogleVerifierOption func(*GoogleVerifier)

// WithGoogleJWKSURL overrides the JWKS endpoint.
func WithGoogleJWKSURL(u string) GoogleVerifierOption {
	return func(v *GoogleVerifier) { v.jwksURL = u }
}

// WithGoogleHTTPClient overrides the HTTP client.
func WithGoogleHTTPClient(c *http.Client) GoogleVerifierOption {
	return func(v *GoogleVerifier) { v.client = c }
}

// WithGoogleClock overrides the clock.
func WithGoogleClock(now func() time.Time) GoogleVerifierOption {
	return func(v *GoogleVerifier) { v.now = now }
}

// NewGoogleVerifier accepts tokens whose audience is any of clientIDs.
func NewGoogleVerifier(clientIDs []string, opts ...GoogleVerifierOption) *GoogleVerifier {
	v := &GoogleVerifier{
		clientIDs: append([]string(nil), clientIDs...),
		jwksURL:   googleDefaultJWKSURL,
		client:    &http.Client{Timeout: 10 * time.Second},
		now:       time.Now,
	}
	for _, o := range opts {
		o(v)
	}
	return v
}

// Verify validates idToken and returns the identity it vouches for.
func (v *GoogleVerifier) Verify(ctx context.Context, idToken string) (*user.ProviderIdentity, error) {
	claims := jwt.MapClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(v.now),
	)
	_, err := parser.ParseWithClaims(idToken, claims, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("missing kid")
		}
		return v.key(ctx, kid)
	})
	if err != nil {
		return nil, fmt.Errorf("google id token rejected: %w", err)
	}

	iss, _ := claims.GetIssuer()
	if iss != "accounts.google.com" && iss != "https://accounts.google.com" {
		return nil, errors.New("google id token rejected: invalid issuer")
	}
	aud, err := claims.GetAudience()
	if err != nil || !v.audienceAllowed(aud) {
		return nil, errors.New("google id token rejected: invalid audience")
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, errors.New("google id token rejected: empty subject")
	}
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)

	verified := false
	switch ev := claims["email_verified"].(type) {
	case bool:
		verified = ev
	case string:
		verified = ev == "true"
	}

	return &user.ProviderIdentity{
		Provider:      user.ProviderGoogle,
		Subject:       sub,
		Email:         email,
		EmailVerified: verified,
		Name:          name,
	}, nil
}

func (v *GoogleVerifier) audienceAllowed(aud []string) bool {
	for _, a := range aud {
		for _, id := range v.clientIDs {
			if id != "" && a == id {
				return true
			}
		}
	}
	return false
}

// key returns the public key for kid. Fresh cached keys are served directly.
// Otherwise (stale cache or unknown kid) the JWKS is refetched, but at most once
// per googleMinRefreshGap whatever the cache state, so an outage of Google's key
// endpoint cannot make every login wait on a fetch. While throttled, or when a
// refetch fails, a previously fetched key for kid is still served. Held under mu
// so concurrent callers share a single fetch.
func (v *GoogleVerifier) key(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	now := v.now()
	known, hasKnown := v.keys[kid]
	if hasKnown && v.fetched && now.Before(v.expiresAt) {
		return known, nil
	}
	if !v.lastFetch.IsZero() && now.Sub(v.lastFetch) < googleMinRefreshGap {
		if hasKnown {
			return known, nil
		}
		return nil, errors.New("unknown kid")
	}
	if err := v.refresh(ctx, now); err != nil {
		if hasKnown {
			return known, nil
		}
		return nil, err
	}
	if k, ok := v.keys[kid]; ok {
		return k, nil
	}
	return nil, errors.New("unknown kid")
}

type googleJWKS struct {
	Keys []struct {
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

func (v *GoogleVerifier) refresh(ctx context.Context, now time.Time) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	if err != nil {
		return err
	}
	// Record the attempt so failures are throttled too.
	v.lastFetch = now
	resp, err := v.client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch google jwks: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch google jwks: status %d", resp.StatusCode)
	}
	var set googleJWKS
	if err := json.NewDecoder(io.LimitReader(resp.Body, providerMaxResponseSize)).Decode(&set); err != nil {
		return fmt.Errorf("decode google jwks: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey, len(set.Keys))
	for _, k := range set.Keys {
		if k.Kty != "RSA" || k.Kid == "" {
			continue
		}
		nb, err1 := base64.RawURLEncoding.DecodeString(k.N)
		eb, err2 := base64.RawURLEncoding.DecodeString(k.E)
		if err1 != nil || err2 != nil || len(nb) == 0 || len(eb) == 0 {
			continue
		}
		e := new(big.Int).SetBytes(eb)
		if !e.IsInt64() || e.Int64() <= 1 || e.Int64() > 1<<31-1 {
			continue
		}
		keys[k.Kid] = &rsa.PublicKey{N: new(big.Int).SetBytes(nb), E: int(e.Int64())}
	}

	ttl := googleDefaultKeyTTL
	if m := googleMaxAgeRe.FindStringSubmatch(resp.Header.Get("Cache-Control")); m != nil {
		if secs, err := strconv.Atoi(m[1]); err == nil && secs >= 0 {
			ttl = time.Duration(secs) * time.Second
		}
	}
	v.keys = keys
	v.expiresAt = now.Add(ttl)
	v.fetched = true
	return nil
}
