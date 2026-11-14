// Package handler exposes the hydration use cases over HTTP.
package handler

import (
	"net/http"

	"fittracker-api/internal/water/service"
	"fittracker-api/pkg/auth"
	"fittracker-api/pkg/httpx"
)

// Router registers the water-intake routes on mux.
func Router(mux *http.ServeMux, signer *auth.Signer, svc *service.Service) {
	mux.HandleFunc(
		"GET /api/profiles/{profileId}/water",
		httpx.RequireAuth(signer, GetDaily(svc)),
	)
	mux.HandleFunc(
		"POST /api/profiles/{profileId}/water",
		httpx.LimitBody(httpx.RequireAuth(signer, Create(svc))),
	)
	mux.HandleFunc(
		"DELETE /api/profiles/{profileId}/water/{id}",
		httpx.RequireAuth(signer, Delete(svc)),
	)
}
