package handler

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/patyukin/go-chat/internal/handler/dto"
	"github.com/patyukin/go-chat/internal/model"
	"sync"
)

type UseCase interface {
	SignUpUseCase(ctx context.Context, in model.SignUpRequest) error
	SignInUseCase(ctx context.Context, in model.SignInRequest) (model.SignInResponse, error)
	GetDomainUseCase() string
	PageMainUseCase(ctx context.Context, id string) (model.MainPageData, error)
	ValidateTokenUseCase(ctx context.Context, token string) (string, error)
	CreateRoomV1UseCase(ctx context.Context, in dto.CreateRoomV1Request) (dto.CreateRoomV1Response, error)
	GetRoomInfoUseCase(ctx context.Context, userUUID, roomUUID string) (dto.PageRoomV1Response, error)
	SaveMassageUseCase(ctx context.Context, roomUUID, userUUID, message string) error
	SetConnectionUseCase(ctx context.Context, connKey string, conn *websocket.Conn) error
	DelConnectionUseCase(ctx context.Context, connKey string) error
	WsHandleUseCase(ctx context.Context, conn *websocket.Conn, roomUUID, userUUID string) error
	GetConnectionKeys(ctx context.Context, roomUUID string) ([]string, error)
	GetConnection(ctx context.Context, connKey string) (*websocket.Conn, error)
	DelConnection(ctx context.Context, connKey string) error
	SelectUserByUUID(ctx context.Context, userUUID string) (model.User, error)
}

type Handler struct {
	uc          UseCase
	connections sync.Map
}

func New(uc UseCase) *Handler {
	return &Handler{uc: uc}
}
