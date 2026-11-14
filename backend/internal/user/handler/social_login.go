package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"fittracker-api/internal/user"
	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/apperror"
	"fittracker-api/pkg/httpx"
)

// googleLoginSchema is the request body of POST /api/auth/oauth/google.
type googleLoginSchema struct {
	IDToken string `json:"idToken"`
}

// githubLoginSchema is the request body of POST /api/auth/oauth/github.
type githubLoginSchema struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"codeVerifier"`
	RedirectURI  string `json:"redirectUri"`
}

// writeSocialBadRequest answers an undecodable or incomplete body with the
// generic social login text, without echoing decoder errors.
func writeSocialBadRequest(w http.ResponseWriter, provider string) {
	apperror.Write(w, apperror.BadRequest(user.ErrSocialLoginFailed(provider).Message))
}

// GoogleLogin handles POST /api/auth/oauth/google.
func GoogleLogin(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body googleLoginSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeSocialBadRequest(w, user.ProviderGoogle)
			return
		}
		idToken := strings.TrimSpace(body.IDToken)
		if idToken == "" {
			writeSocialBadRequest(w, user.ProviderGoogle)
			return
		}

		result, err := svc.LoginWithGoogle(r.Context(), idToken)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}

// GitHubLogin handles POST /api/auth/oauth/github.
func GitHubLogin(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body githubLoginSchema
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeSocialBadRequest(w, user.ProviderGitHub)
			return
		}
		code := strings.TrimSpace(body.Code)
		verifier := strings.TrimSpace(body.CodeVerifier)
		redirectURI := strings.TrimSpace(body.RedirectURI)
		if code == "" || verifier == "" || redirectURI == "" {
			writeSocialBadRequest(w, user.ProviderGitHub)
			return
		}

		result, err := svc.LoginWithGitHub(r.Context(), code, verifier, redirectURI)
		if err != nil {
			apperror.Write(w, err)
			return
		}

		httpx.JSONResponse(w, http.StatusOK, result)
	}
}
