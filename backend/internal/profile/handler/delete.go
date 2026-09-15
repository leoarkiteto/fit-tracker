package handler

import (
	"net/http"

	"fittracker-api/internal/profile/service"
	"fittracker-api/pkg/apperror"
)

// Delete handles DELETE /api/profiles/{id}.
func Delete(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.Delete(r.Context(), r.PathValue("id")); err != nil {
			apperror.Write(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
