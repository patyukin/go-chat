package handler

import (
	"github.com/patyukin/go-chat/internal/model"
	"github.com/patyukin/go-chat/pkg/httperror"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"time"
)

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error().Err(err).Msgf("failed to sign in, error: %v", err)
		httperror.SendError(w, "invalid sign in: r.ParseForm()", http.StatusBadRequest)
		return
	}

	m := model.SignInRequest{Login: r.FormValue("login"), Password: r.FormValue("password")}
	tokens, err := h.uc.SignInUseCase(r.Context(), m)
	if err != nil {
		log.Error().Err(err).Msgf("failed to sign in, error: %v", err)
		httperror.SendError(w, "invalid sign in: h.uc.SignInUseCase", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Domain:   h.uc.GetDomainUseCase(),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Domain:   h.uc.GetDomainUseCase(),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) PageSignInHandler(w http.ResponseWriter, r *http.Request) {
	signInTemplate, err := template.ParseFiles("internal/templates/sign-in.html")
	err = signInTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}
