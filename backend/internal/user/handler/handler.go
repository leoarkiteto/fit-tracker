// Package handler exposes the account use cases over HTTP.
package handler

import (
	"net/http"

	"fittracker-api/internal/user/service"
	"fittracker-api/pkg/auth"
	"fittracker-api/pkg/httpx"
)

// Router registers the account routes on mux.
func Router(mux *http.ServeMux, signer *auth.Signer, svc *service.Service) {
	mux.HandleFunc("POST /api/auth/register", httpx.LimitBody(Register(svc)))
	mux.HandleFunc("POST /api/auth/login", httpx.LimitBody(Login(svc)))
	mux.HandleFunc("GET /api/auth/me", httpx.RequireAuth(signer, Me(svc)))
	mux.HandleFunc(
		"PATCH /api/auth/change-password",
		httpx.LimitBody(httpx.RequireAuth(signer, ChangePassword(svc))),
	)
}
