package infrastructure

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"fittracker-api/internal/user"
)

const googleTestClientID = "client-1.apps.googleusercontent.com"

type googleFixture struct {
	t      *testing.T
	server *httptest.Server
	hits   atomic.Int32
	down   atomic.Bool // when set, the JWKS endpoint answers 503
	mu     sync.Mutex
	keys   map[string]*rsa.PublicKey
	cc     string
	now    time.Time
}

func googleNewKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func googleNewFixture(t *testing.T, cacheControl string, keys map[string]*rsa.PublicKey) *googleFixture {
	t.Helper()
	f := &googleFixture{t: t, keys: keys, cc: cacheControl, now: time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.hits.Add(1)
		if f.down.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		type jwk struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			N   string `json:"n"`
			E   string `json:"e"`
		}
		var set struct {
			Keys []jwk `json:"keys"`
		}
		for kid, pub := range f.keys {
			set.Keys = append(set.Keys, jwk{
				Kid: kid, Kty: "RSA",
				N: base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
				E: base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
			})
		}
		if f.cc != "" {
			w.Header().Set("Cache-Control", f.cc)
		}
		_ = json.NewEncoder(w).Encode(set)
	}))
	t.Cleanup(f.server.Close)
	return f
}

func (f *googleFixture) verifier(clientIDs ...string) *GoogleVerifier {
	if len(clientIDs) == 0 {
		clientIDs = []string{googleTestClientID}
	}
	return NewGoogleVerifier(clientIDs,
		WithGoogleJWKSURL(f.server.URL),
		WithGoogleClock(func() time.Time { return f.now }),
	)
}

func (f *googleFixture) claims(mut func(jwt.MapClaims)) jwt.MapClaims {
	c := jwt.MapClaims{
		"iss":            "https://accounts.google.com",
		"aud":            googleTestClientID,
		"sub":            "g-123",
		"email":          "ana@example.com",
		"email_verified": true,
		"name":           "Ana",
		"iat":            f.now.Add(-time.Minute).Unix(),
		"exp":            f.now.Add(time.Hour).Unix(),
	}
	if mut != nil {
		mut(c)
	}
	return c
}

func googleSign(t *testing.T, key *rsa.PrivateKey, kid string, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = kid
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestGoogleVerifier_Valid(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "public, max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	tok := googleSign(t, key, "k1", f.claims(nil))

	id, err := f.verifier().Verify(context.Background(), tok)
	if err != nil {
		t.Fatal(err)
	}
	want := user.ProviderIdentity{Provider: user.ProviderGoogle, Subject: "g-123", Email: "ana@example.com", EmailVerified: true, Name: "Ana"}
	if *id != want {
		t.Fatalf("got %+v want %+v", *id, want)
	}
}

func TestGoogleVerifier_Rejections(t *testing.T) {
	key := googleNewKey(t)
	other := googleNewKey(t)
	f := googleNewFixture(t, "max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})

	cases := map[string]string{
		"bad signature": googleSign(t, other, "k1", f.claims(nil)),
		"wrong aud":     googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { c["aud"] = "someone-else" })),
		"wrong iss":     googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { c["iss"] = "https://evil.example.com" })),
		"expired":       googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { c["exp"] = f.now.Add(-time.Minute).Unix() })),
		"no exp":        googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { delete(c, "exp") })),
		"unknown kid":   googleSign(t, key, "nope", f.claims(nil)),
		"empty sub":     googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { c["sub"] = "" })),
		"garbage":       "not-a-jwt",
	}
	for name, tok := range cases {
		t.Run(name, func(t *testing.T) {
			if id, err := f.verifier().Verify(context.Background(), tok); err == nil {
				t.Fatalf("expected error, got %+v", id)
			}
		})
	}
}

func TestGoogleVerifier_RejectsNonRS256(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	tok := jwt.NewWithClaims(jwt.SigningMethodRS384, f.claims(nil))
	tok.Header["kid"] = "k1"
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.verifier().Verify(context.Background(), s); err == nil {
		t.Fatal("expected RS384 to be rejected")
	}
}

func TestGoogleVerifier_BothIssuers(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	tok := googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { c["iss"] = "accounts.google.com" }))
	if _, err := f.verifier().Verify(context.Background(), tok); err != nil {
		t.Fatal(err)
	}
}

func TestGoogleVerifier_MultiClientAudience(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	tok := googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { c["aud"] = "second-id" }))
	if _, err := f.verifier("first-id", "second-id").Verify(context.Background(), tok); err != nil {
		t.Fatal(err)
	}
	if _, err := f.verifier("first-id", "third-id").Verify(context.Background(), tok); err == nil {
		t.Fatal("expected audience mismatch")
	}
}

func TestGoogleVerifier_EmailNotVerified(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	for name, val := range map[string]any{"bool false": false, "string false": "false"} {
		t.Run(name, func(t *testing.T) {
			tok := googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { c["email_verified"] = val }))
			id, err := f.verifier().Verify(context.Background(), tok)
			if err != nil {
				t.Fatal(err)
			}
			if id.EmailVerified {
				t.Fatal("EmailVerified must be false")
			}
		})
	}
}

func TestGoogleVerifier_EmailVerifiedString(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	tok := googleSign(t, key, "k1", f.claims(func(c jwt.MapClaims) { c["email_verified"] = "true" }))
	id, err := f.verifier().Verify(context.Background(), tok)
	if err != nil {
		t.Fatal(err)
	}
	if !id.EmailVerified {
		t.Fatal("string \"true\" must count as verified")
	}
}

func TestGoogleVerifier_CachesUntilMaxAge(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "public, max-age=300", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	v := f.verifier()
	verify := func() {
		t.Helper()
		tok := googleSign(t, key, "k1", f.claims(nil))
		if _, err := v.Verify(context.Background(), tok); err != nil {
			t.Fatal(err)
		}
	}
	verify()
	verify()
	if n := f.hits.Load(); n != 1 {
		t.Fatalf("hits = %d, want 1", n)
	}
	f.now = f.now.Add(299 * time.Second)
	verify()
	if n := f.hits.Load(); n != 1 {
		t.Fatalf("hits before expiry = %d, want 1", n)
	}
	f.now = f.now.Add(2 * time.Second)
	verify()
	if n := f.hits.Load(); n != 2 {
		t.Fatalf("hits after expiry = %d, want 2", n)
	}
}

func TestGoogleVerifier_DefaultTTLWithoutCacheControl(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	v := f.verifier()
	for i := 0; i < 2; i++ {
		if _, err := v.Verify(context.Background(), googleSign(t, key, "k1", f.claims(nil))); err != nil {
			t.Fatal(err)
		}
	}
	if n := f.hits.Load(); n != 1 {
		t.Fatalf("hits = %d, want 1", n)
	}
}

func TestGoogleVerifier_UnknownKidRefreshThrottle(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=86400", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	v := f.verifier()
	ctx := context.Background()

	if _, err := v.Verify(ctx, googleSign(t, key, "k1", f.claims(nil))); err != nil {
		t.Fatal(err)
	}
	if n := f.hits.Load(); n != 1 {
		t.Fatalf("initial hits = %d", n)
	}

	// Rotation: new key appears upstream, but within a minute of the last fetch we do not refetch.
	newKey := googleNewKey(t)
	f.mu.Lock()
	f.keys["k2"] = &newKey.PublicKey
	f.mu.Unlock()
	f.now = f.now.Add(10 * time.Second)
	tok2 := googleSign(t, newKey, "k2", f.claims(nil))
	for i := 0; i < 5; i++ {
		if _, err := v.Verify(ctx, tok2); err == nil {
			t.Fatal("expected unknown kid error while throttled")
		}
	}
	if n := f.hits.Load(); n != 1 {
		t.Fatalf("hits while throttled = %d, want 1", n)
	}

	// After a minute one forced refresh happens and picks up the new key.
	f.now = f.now.Add(time.Minute)
	if _, err := v.Verify(ctx, tok2); err != nil {
		t.Fatal(err)
	}
	if n := f.hits.Load(); n != 2 {
		t.Fatalf("hits after minute = %d, want 2", n)
	}

	// Still-unknown kid: one refetch, then throttled again.
	bogus := googleSign(t, key, "bogus", f.claims(nil))
	f.now = f.now.Add(2 * time.Minute)
	for i := 0; i < 3; i++ {
		if _, err := v.Verify(ctx, bogus); err == nil {
			t.Fatal("expected error")
		}
	}
	if n := f.hits.Load(); n != 3 {
		t.Fatalf("hits = %d, want 3", n)
	}
}

// After the cache expires, a JWKS outage must not make every login fetch:
// refetches stay at most one per minute and the last known key keeps working.
func TestGoogleVerifier_StaleCacheOutageThrottledWithFallback(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=60", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	v := f.verifier()
	ctx := context.Background()
	verifyOK := func() {
		t.Helper()
		if _, err := v.Verify(ctx, googleSign(t, key, "k1", f.claims(nil))); err != nil {
			t.Fatalf("known key must keep working during an outage: %v", err)
		}
	}

	verifyOK()
	if n := f.hits.Load(); n != 1 {
		t.Fatalf("initial hits = %d, want 1", n)
	}

	f.down.Store(true)
	f.now = f.now.Add(2 * time.Minute) // cache expired
	verifyOK()                         // one refetch attempt fails, stale key served
	if n := f.hits.Load(); n != 2 {
		t.Fatalf("hits after expiry = %d, want 2", n)
	}

	f.now = f.now.Add(10 * time.Second)
	for i := 0; i < 5; i++ {
		verifyOK()
	}
	if n := f.hits.Load(); n != 2 {
		t.Fatalf("hits while throttled during outage = %d, want 2", n)
	}

	f.now = f.now.Add(time.Minute)
	verifyOK()
	if n := f.hits.Load(); n != 3 {
		t.Fatalf("hits after a minute = %d, want 3", n)
	}

	f.down.Store(false)
	f.now = f.now.Add(time.Minute)
	verifyOK()
	if n := f.hits.Load(); n != 4 {
		t.Fatalf("hits after recovery = %d, want 4", n)
	}
	f.now = f.now.Add(10 * time.Second)
	verifyOK()
	if n := f.hits.Load(); n != 4 {
		t.Fatalf("fresh cache after recovery must not refetch: hits = %d, want 4", n)
	}
}

// When the very first fetch fails there is no key to fall back on; failed
// attempts are still throttled.
func TestGoogleVerifier_InitialOutageThrottled(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	f.down.Store(true)
	v := f.verifier()
	tok := googleSign(t, key, "k1", f.claims(nil))

	for i := 0; i < 3; i++ {
		if _, err := v.Verify(context.Background(), tok); err == nil {
			t.Fatal("expected error while JWKS is unavailable")
		}
	}
	if n := f.hits.Load(); n != 1 {
		t.Fatalf("hits = %d, want 1", n)
	}
}

func TestGoogleVerifier_ConcurrentUnknownKidSingleFetch(t *testing.T) {
	key := googleNewKey(t)
	f := googleNewFixture(t, "max-age=3600", map[string]*rsa.PublicKey{"k1": &key.PublicKey})
	v := f.verifier()
	tok := googleSign(t, key, "k1", f.claims(nil))

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := v.Verify(context.Background(), tok); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if n := f.hits.Load(); n != 1 {
		t.Fatalf("hits = %d, want 1", n)
	}
}
