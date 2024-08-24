package handler

import (
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

func (h *Handler) PageMainHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Header.Get("X-User-UUID")
	if userUUID == "" {
		log.Error().Msgf("unable to get user uuid")
		http.Error(w, "Unable to render page", http.StatusBadRequest)
		return
	}

	data, err := h.uc.PageMainUseCase(r.Context(), userUUID)
	if err != nil {
		log.Error().Msgf("failed to get main page data, error: %v", err)
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	log.Debug().Msgf("data: %+v", data)

	mainTemplate, err := template.ParseFiles("internal/templates/main.html")
	if err = mainTemplate.Execute(w, data); err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}
