package usecase

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
)

func (u *UseCase) SetConnectionUseCase(ctx context.Context, connKey string, conn *websocket.Conn) error {
	err := u.chr.SetConnection(ctx, connKey, conn)
	if err != nil {
		return fmt.Errorf("unable to set connection: %w", err)
	}

	return nil
}

func (u *UseCase) DelConnectionUseCase(ctx context.Context, connKey string) error {
	err := u.chr.DelConnection(ctx, connKey)
	if err != nil {
		return fmt.Errorf("unable to del connection: %w", err)
	}

	return nil
}
