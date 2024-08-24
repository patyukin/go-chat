package router

import (
	"github.com/patyukin/go-chat/internal/handler"
	"github.com/patyukin/go-chat/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Handler interface {
	SignUpHandler(w http.ResponseWriter, r *http.Request)
	SignInHandler(w http.ResponseWriter, r *http.Request)
	PageSignUpHandler(w http.ResponseWriter, r *http.Request)
	PageSignInHandler(w http.ResponseWriter, r *http.Request)
	PageMainHandler(w http.ResponseWriter, r *http.Request)
	CreateRoomV1Handler(w http.ResponseWriter, r *http.Request)
}

func Init(h *handler.Handler) http.Handler {
	r := http.ServeMux{}

	prometheus.MustRegister(metrics.IncomingTraffic)

	r.Handle("POST /sign-up", http.HandlerFunc(h.SignUpHandler))
	r.Handle("GET /sign-up", http.HandlerFunc(h.PageSignUpHandler))

	r.Handle("POST /sign-in", http.HandlerFunc(h.SignInHandler))
	r.Handle("GET /sign-in", http.HandlerFunc(h.PageSignInHandler))

	r.Handle("GET /", h.AuthMiddleware(http.HandlerFunc(h.PageMainHandler)))
	r.Handle("GET /rooms/{room_id}", h.AuthMiddleware(http.HandlerFunc(h.PageRoomHandler)))
	r.Handle("POST /v1/create-room", h.AuthMiddleware(http.HandlerFunc(h.CreateRoomV1Handler)))
	r.Handle("GET /v1/rooms/{room_id}/messages", h.AuthMiddleware(http.HandlerFunc(h.CreateRoomV1Handler)))
	r.Handle("GET /v1/rooms/{room_id}/users", h.AuthMiddleware(http.HandlerFunc(h.CreateRoomV1Handler)))

	r.Handle("GET /ws/rooms/", h.AuthMiddleware(http.HandlerFunc(h.WsRoomHandler)))

	r.Handle("GET /metrics", promhttp.Handler())

	return &r
}
