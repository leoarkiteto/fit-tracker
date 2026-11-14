// Package handler exposes the bioimpedance use cases over HTTP.
package handler

import (
	"net/http"

	"fittracker-api/internal/bioimpedance/service"
	"fittracker-api/pkg/auth"
	"fittracker-api/pkg/httpx"
)

// Router registers the bioimpedance routes on mux.
func Router(mux *http.ServeMux, signer *auth.Signer, svc *service.Service) {
	mux.HandleFunc(
		"GET /api/profiles/{profileId}/bioimpedance",
		httpx.RequireAuth(signer, GetAll(svc)),
	)
	mux.HandleFunc(
		"POST /api/profiles/{profileId}/bioimpedance",
		httpx.LimitBody(httpx.RequireAuth(signer, Create(svc))),
	)
	mux.HandleFunc(
		"GET /api/profiles/{profileId}/bioimpedance/latest",
		httpx.RequireAuth(signer, GetLatest(svc)),
	)
	mux.HandleFunc(
		"DELETE /api/profiles/{profileId}/bioimpedance/{id}",
		httpx.RequireAuth(signer, Delete(svc)),
	)
}
