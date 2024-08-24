package handler

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) PageRoomStartHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Header.Get("X-User-UUID")
	if userUUID == "" {
		log.Error().Msgf("unable to get user uuid")
		http.Error(w, "Unable to render page", http.StatusBadRequest)
		return
	}
}
