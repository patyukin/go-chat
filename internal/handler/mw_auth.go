package handler

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(accessToken)
		if err != nil {
			log.Error().Msgf("failed to get cookie, error: %v", err)
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}

		log.Debug().Msgf("cookie name: %s, value: %s", cookie.Name, cookie.Value)

		userUUID, err := h.uc.ValidateTokenUseCase(r.Context(), cookie.Value)
		if err != nil {
			log.Error().Msgf("failed to validate token, error: %v", err)
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}

		log.Debug().Msgf("user uuid: %v", userUUID)

		r.Header.Set("X-User-UUID", userUUID)

		next.ServeHTTP(w, r)
	})
}
