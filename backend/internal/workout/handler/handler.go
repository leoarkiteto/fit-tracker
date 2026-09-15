// Package handler exposes the workout use cases over HTTP.
package handler

import (
	"net/http"

	"fittracker-api/internal/workout/service"
	"fittracker-api/pkg/auth"
	"fittracker-api/pkg/httpx"
)

// Router registers the workout routes on mux.
func Router(mux *http.ServeMux, signer *auth.Signer, svc *service.Service) {
	mux.HandleFunc(
		"GET /api/profiles/{profileId}/workouts",
		httpx.RequireAuth(signer, GetAll(svc)),
	)
	mux.HandleFunc(
		"GET /api/profiles/{profileId}/workouts/today",
		httpx.RequireAuth(signer, GetToday(svc)),
	)
	mux.HandleFunc(
		"POST /api/profiles/{profileId}/workouts",
		httpx.LimitBody(httpx.RequireAuth(signer, Create(svc))),
	)
	mux.HandleFunc(
		"GET /api/profiles/{profileId}/workouts/{id}",
		httpx.RequireAuth(signer, GetByID(svc)),
	)
	mux.HandleFunc(
		"PUT /api/profiles/{profileId}/workouts/{id}",
		httpx.LimitBody(httpx.RequireAuth(signer, Update(svc))),
	)
	mux.HandleFunc(
		"DELETE /api/profiles/{profileId}/workouts/{id}",
		httpx.RequireAuth(signer, Delete(svc)),
	)
}
