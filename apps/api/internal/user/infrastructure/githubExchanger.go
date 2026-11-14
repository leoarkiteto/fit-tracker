package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fittracker-api/internal/user"
)

// GitHubExchanger trades a GitHub authorization code for the user's identity.
type GitHubExchanger struct {
	clientID     string
	clientSecret string
	webBase      string
	apiBase      string
	client       *http.Client
}

var _ user.GitHubExchanger = (*GitHubExchanger)(nil)

// GitHubExchangerOption customises a GitHubExchanger.
type GitHubExchangerOption func(*GitHubExchanger)

// WithGitHubBaseURLs overrides the OAuth host and the API host.
func WithGitHubBaseURLs(web, api string) GitHubExchangerOption {
	return func(g *GitHubExchanger) {
		g.webBase = strings.TrimRight(web, "/")
		g.apiBase = strings.TrimRight(api, "/")
	}
}

// WithGitHubHTTPClient overrides the HTTP client.
func WithGitHubHTTPClient(c *http.Client) GitHubExchangerOption {
	return func(g *GitHubExchanger) { g.client = c }
}

// NewGitHubExchanger builds an exchanger for the given OAuth app credentials.
func NewGitHubExchanger(clientID, clientSecret string, opts ...GitHubExchangerOption) *GitHubExchanger {
	g := &GitHubExchanger{
		clientID:     clientID,
		clientSecret: clientSecret,
		webBase:      "https://github.com",
		apiBase:      "https://api.github.com",
		client:       &http.Client{Timeout: 10 * time.Second},
	}
	for _, o := range opts {
		o(g)
	}
	return g
}

// Exchange trades the code for an access token (never stored or logged) and reads the profile.
func (g *GitHubExchanger) Exchange(ctx context.Context, code, codeVerifier, redirectURI string) (*user.ProviderIdentity, error) {
	token, err := g.accessToken(ctx, code, codeVerifier, redirectURI)
	if err != nil {
		return nil, err
	}

	var profile struct {
		ID    int64   `json:"id"`
		Login string  `json:"login"`
		Name  *string `json:"name"`
	}
	if err := g.getJSON(ctx, "/user", token, &profile); err != nil {
		return nil, err
	}
	if profile.ID == 0 {
		return nil, errors.New("github user has no id")
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := g.getJSON(ctx, "/user/emails", token, &emails); err != nil {
		return nil, err
	}

	id := &user.ProviderIdentity{
		Provider: user.ProviderGitHub,
		Subject:  strconv.FormatInt(profile.ID, 10),
		Login:    profile.Login,
	}
	if profile.Name != nil {
		id.Name = *profile.Name
	}
	for _, e := range emails {
		if e.Primary && e.Verified && e.Email != "" {
			id.Email = e.Email
			id.EmailVerified = true
			break
		}
	}
	return id, nil
}

func (g *GitHubExchanger) accessToken(ctx context.Context, code, codeVerifier, redirectURI string) (string, error) {
	form := url.Values{
		"client_id":     {g.clientID},
		"client_secret": {g.clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"code_verifier": {codeVerifier},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.webBase+"/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", errors.New("build github token request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", errors.New("github token exchange failed: request error")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("github token exchange failed: status %d", resp.StatusCode)
	}
	var body struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, providerMaxResponseSize)).Decode(&body); err != nil {
		return "", errors.New("github token exchange failed: invalid response")
	}
	if body.Error != "" {
		return "", fmt.Errorf("github token exchange failed: %s", body.Error)
	}
	if body.AccessToken == "" {
		return "", errors.New("github token exchange failed: no access token")
	}
	return body.AccessToken, nil
}

func (g *GitHubExchanger) getJSON(ctx context.Context, path, token string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.apiBase+path, nil)
	if err != nil {
		return errors.New("build github request")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("github %s failed: request error", path)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("github %s failed: status %d", path, resp.StatusCode)
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, providerMaxResponseSize)).Decode(out); err != nil {
		return fmt.Errorf("github %s failed: invalid response", path)
	}
	return nil
}
