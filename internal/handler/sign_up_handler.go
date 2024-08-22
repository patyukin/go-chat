package handler

import (
	"github.com/patyukin/go-chat/internal/model"
	"github.com/patyukin/go-chat/pkg/httperror"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	m := model.SignUpRequest{Login: r.FormValue("login"), Password: r.FormValue("password")}
	log.Debug().Msgf("sign up request: %+v", m)

	if err := h.uc.SignUpUseCase(r.Context(), m); err != nil {
		log.Error().Err(err).Msgf("failed to sign up, error: %v", err)
		httperror.SendError(w, "invalid sign up", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/sign-in", http.StatusFound)
}

func (h *Handler) PageSignUpHandler(w http.ResponseWriter, r *http.Request) {
	signInTemplate, err := template.ParseFiles("internal/templates/sign-up.html")
	err = signInTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}
