package usecase

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/patyukin/go-chat/internal/auth"
	"github.com/patyukin/go-chat/internal/cacher"
	"github.com/patyukin/go-chat/internal/config"
	"github.com/patyukin/go-chat/internal/db"
	"github.com/patyukin/go-chat/internal/model"
)

type UseCase struct {
	cfg        *config.Config
	registry   *db.Client
	authClient *auth.Client
	chr        *cacher.Cacher
}

func New(cfg *config.Config, registry *db.Client, authClient *auth.Client, chr *cacher.Cacher) *UseCase {
	return &UseCase{
		cfg:        cfg,
		registry:   registry,
		authClient: authClient,
		chr:        chr,
	}
}

func (u *UseCase) GetDomainUseCase() string {
	return u.cfg.CookieDomain
}

func (u *UseCase) GetConnectionKeys(ctx context.Context, roomUUID string) ([]string, error) {
	keys, err := u.chr.GetConnectionKeys(ctx, roomUUID)
	if err != nil {
		return nil, fmt.Errorf("unable to get connection keys: %w", err)
	}

	return keys, nil
}

func (u *UseCase) GetConnection(ctx context.Context, connKey string) (*websocket.Conn, error) {
	conn, err := u.chr.GetConnection(ctx, connKey)
	if err != nil {
		return nil, fmt.Errorf("unable to get connection: %w", err)
	}

	return conn, nil
}

func (u *UseCase) DelConnection(ctx context.Context, connKey string) error {
	err := u.chr.DelConnection(ctx, connKey)
	if err != nil {
		return fmt.Errorf("unable to del connection: %w", err)
	}

	return nil
}

func (u *UseCase) SelectUserByUUID(ctx context.Context, userUUID string) (model.User, error) {
	user, err := u.registry.GetRepo().SelectUserByUUID(ctx, userUUID)
	if err != nil {
		return model.User{}, fmt.Errorf("unable to get user: %w", err)
	}

	return user, nil
}
