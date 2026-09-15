package handler

import (
	"net/http"

	"fittracker-api/internal/bioimpedance/service"
	"fittracker-api/pkg/apperror"
)

// Delete handles DELETE /api/profiles/{profileId}/bioimpedance/{id}.
func Delete(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := svc.Delete(r.Context(), r.PathValue("id"), r.PathValue("profileId"))
		if err != nil {
			apperror.Write(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
