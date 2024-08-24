package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/go-chat/internal/db"
	"github.com/patyukin/go-chat/internal/handler/dto"
	"github.com/patyukin/go-chat/internal/model"
)

func (u *UseCase) GetRoomInfoUseCase(ctx context.Context, userUUID, roomUUID string) (dto.PageRoomV1Response, error) {
	var room model.Room
	var user model.User
	var messages []model.Message
	var users []model.User
	var err error

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		room, err = repo.SelectRoomByRoomUUID(ctx, roomUUID)
		if err != nil {
			return fmt.Errorf("failed SelectRoomByRoomUUID: %w", err)
		}

		messages, err = repo.SelectMessagesByRoomUUIDUserUUID(ctx, userUUID, roomUUID)
		if err != nil {
			return fmt.Errorf("failed SelectMessagesByRoomUUIDUserUUID: %w", err)
		}

		users, err = repo.SelectAllUsers(ctx)
		if err != nil {
			return fmt.Errorf("failed SelectUsersWithoutUserUUID: %w", err)
		}

		user, err = repo.SelectUserByUUID(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed SelectUserByUUID: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.PageRoomV1Response{}, fmt.Errorf("failed ReadCommitted: %w", err)
	}

	return dto.PageRoomV1Response{Messages: messages, Users: users, Room: room, User: user}, nil
}
