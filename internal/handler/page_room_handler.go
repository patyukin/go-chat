package handler

import (
	"github.com/patyukin/go-chat/pkg/httperror"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

func (h *Handler) PageRoomHandler(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Header.Get("X-User-UUID")
	if userUUID == "" {
		log.Error().Msgf("unable to get user uuid")
		httperror.SendError(w, "Unable to render page", http.StatusBadRequest)
		return
	}

	roomUUID := r.PathValue("room_id")
	if roomUUID == "" {
		log.Error().Msgf("unable to get room room_id")
		httperror.SendError(w, "Unable to render page", http.StatusBadRequest)
		return
	}

	pageData, err := h.uc.GetRoomInfoUseCase(r.Context(), userUUID, roomUUID)

	log.Debug().Msgf("pageData: %+v", pageData)

	if err != nil {
		log.Error().Msgf("failed to get room info, error: %v", err)
		httperror.SendError(w, "Unable to render page", http.StatusBadRequest)
		return
	}

	signInTemplate, err := template.ParseFiles("internal/templates/room.html")
	if err != nil {
		log.Error().Msgf("failed to render page, error: %v", err)
		httperror.SendError(w, "Unable to render page", http.StatusInternalServerError)
		return
	}

	if err = signInTemplate.Execute(w, pageData); err != nil {
		log.Error().Msgf("failed to render page, error: %v", err)
		httperror.SendError(w, "Unable to render page", http.StatusInternalServerError)
		return
	}
}
