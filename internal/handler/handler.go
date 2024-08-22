package handler

import (
	"context"
	"github.com/patyukin/go-chat/internal/model"
	"net/http"
)

type UseCase interface {
	SignUpUseCase(ctx context.Context, in model.SignUpRequest) error
	SignInUseCase(ctx context.Context, in model.SignInRequest) (model.SignInResponse, error)
	GetDomainUseCase() string
}

type Handler struct {
	uc UseCase
}

func New(uc UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) PageMainHandler(w http.ResponseWriter, r *http.Request) {

}
