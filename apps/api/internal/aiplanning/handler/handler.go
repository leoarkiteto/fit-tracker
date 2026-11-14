// Package handler exposes the AI-planning use cases over HTTP.
package handler

import (
	"net/http"

	"fittracker-api/internal/aiplanning/service"
	"fittracker-api/pkg/auth"
	"fittracker-api/pkg/httpx"
)

// Router registers the AI-planning routes on mux.
func Router(mux *http.ServeMux, signer *auth.Signer, svc *service.Service) {
	mux.HandleFunc(
		"POST /api/ai/planning/generate",
		httpx.LimitBody(httpx.RequireAuth(signer, Generate(svc))),
	)
	mux.HandleFunc(
		"POST /api/ai/planning/accept",
		httpx.LimitBody(httpx.RequireAuth(signer, Accept(svc))),
	)
	mux.HandleFunc("GET /api/ai/planning/status", Status(svc))
}
