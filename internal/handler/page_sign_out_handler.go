package handler

import (
	"net/http"
	"time"
)

func (h *Handler) PageSignOutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessToken,
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		Path:     "/",
		Domain:   h.uc.GetDomainUseCase(),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     accessToken,
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		Path:     "/",
		Domain:   h.uc.GetDomainUseCase(),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/sign-in", http.StatusFound)
}
