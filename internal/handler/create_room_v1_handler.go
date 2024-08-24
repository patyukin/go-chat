package handler

import (
	"encoding/json"
	"github.com/patyukin/go-chat/internal/handler/dto"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) CreateRoomV1Handler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("CreateRoomV1Handler")
	userUUID := r.Header.Get("X-User-UUID")
	if userUUID == "" {
		log.Error().Msgf("unable to get user uuid")
		http.Error(w, "Unable to render page", http.StatusBadRequest)
		return
	}

	var in dto.CreateRoomV1Request
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		log.Error().Msgf("unable to parse form: %v", err)
		http.Error(w, "Unable to render page", http.StatusBadRequest)
		return
	}

	in.Users = append(in.Users, userUUID)

	out, err := h.uc.CreateRoomV1UseCase(r.Context(), in)
	if err != nil {
		log.Error().Msgf("failed to create room, error: %v", err)
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		return
	}

	log.Debug().Msgf("out: %+v", out)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		return
	}
}
