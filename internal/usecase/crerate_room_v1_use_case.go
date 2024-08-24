package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/go-chat/internal/db"
	"github.com/patyukin/go-chat/internal/handler/dto"
)

func (u *UseCase) CreateRoomV1UseCase(ctx context.Context, in dto.CreateRoomV1Request) (dto.CreateRoomV1Response, error) {
	var roomID string
	var err error

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		roomID, err = repo.InsertIntoRooms(ctx, in.Name)
		if err != nil {
			return fmt.Errorf("failed InsertIntoRooms from CreateRoomV1UseCase: %w", err)
		}

		err = repo.InsertIntoRoomsUsers(ctx, roomID, in.Users)
		if err != nil {
			return fmt.Errorf("failed InsertIntoRoomsUsers from CreateRoomV1UseCase: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.CreateRoomV1Response{}, fmt.Errorf("failed ReadCommitted from CreateRoomV1UseCase: %w", err)
	}

	return dto.CreateRoomV1Response{RoomID: roomID}, nil
}
