// Package handler exposes the completed-workout use cases over HTTP.
package handler

import (
	"net/http"

	"fittracker-api/internal/completedworkout/service"
	"fittracker-api/pkg/auth"
	"fittracker-api/pkg/httpx"
)

// Router registers the training-history routes on mux.
func Router(mux *http.ServeMux, signer *auth.Signer, svc *service.Service) {
	mux.HandleFunc(
		"GET /api/profiles/{profileId}/completed-workouts",
		httpx.RequireAuth(signer, GetAll(svc)),
	)
	mux.HandleFunc(
		"POST /api/profiles/{profileId}/completed-workouts",
		httpx.LimitBody(httpx.RequireAuth(signer, Complete(svc))),
	)
	mux.HandleFunc(
		"GET /api/profiles/{profileId}/completed-workouts/stats",
		httpx.RequireAuth(signer, GetStats(svc)),
	)
}
