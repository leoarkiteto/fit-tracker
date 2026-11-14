package mock

import (
	"context"

	"fittracker-api/internal/user"
)

// GoogleVerifier is a fake user.GoogleVerifier returning Identity or Err.
type GoogleVerifier struct {
	Err      error
	Identity *user.ProviderIdentity

	Calls     int
	LastToken string
}

var _ user.GoogleVerifier = (*GoogleVerifier)(nil)

// Verify records the token and returns the configured result.
func (v *GoogleVerifier) Verify(_ context.Context, idToken string) (*user.ProviderIdentity, error) {
	v.Calls++
	v.LastToken = idToken
	if v.Err != nil {
		return nil, v.Err
	}
	return v.Identity, nil
}

// GitHubExchanger is a fake user.GitHubExchanger returning Identity or Err.
type GitHubExchanger struct {
	Err      error
	Identity *user.ProviderIdentity

	Calls           int
	LastCode        string
	LastVerifier    string
	LastRedirectURI string
}

var _ user.GitHubExchanger = (*GitHubExchanger)(nil)

// Exchange records the arguments and returns the configured result.
func (e *GitHubExchanger) Exchange(_ context.Context, code, codeVerifier, redirectURI string) (*user.ProviderIdentity, error) {
	e.Calls++
	e.LastCode = code
	e.LastVerifier = codeVerifier
	e.LastRedirectURI = redirectURI
	if e.Err != nil {
		return nil, e.Err
	}
	return e.Identity, nil
}
