// Package handler exposes the profile use cases over HTTP.
package handler

import (
	"net/http"

	"fittracker-api/internal/profile/service"
	"fittracker-api/pkg/auth"
	"fittracker-api/pkg/httpx"
)

// Router registers the profile routes on mux.
func Router(mux *http.ServeMux, signer *auth.Signer, svc *service.Service) {
	mux.HandleFunc("GET /api/profiles", httpx.RequireAuth(signer, GetAll(svc)))
	mux.HandleFunc("GET /api/profiles/{id}", httpx.RequireAuth(signer, GetByID(svc)))
	mux.HandleFunc(
		"PUT /api/profiles/{id}",
		httpx.LimitBody(httpx.RequireAuth(signer, Update(svc))),
	)
	mux.HandleFunc("DELETE /api/profiles/{id}", httpx.RequireAuth(signer, Delete(svc)))
}
