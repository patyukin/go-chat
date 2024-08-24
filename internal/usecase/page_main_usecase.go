package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/go-chat/internal/db"
	"github.com/patyukin/go-chat/internal/model"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) PageMainUseCase(ctx context.Context, id string) (model.MainPageData, error) {
	var users []model.User
	var rooms []model.Room
	var user model.User
	var err error

	err = u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		users, err = repo.SelectUsersWithoutUserUUID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed SelectUsersWithoutUserUUID: %w", err)
		}

		log.Debug().Msgf("users: %v", users)

		rooms, err = repo.SelectAllRooms(ctx, id)
		if err != nil {
			return fmt.Errorf("failed SelectAllRooms: %w", err)
		}

		log.Debug().Msgf("rooms: %v", rooms)

		user, err = repo.SelectUserByUUID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed SelectUserByUUID: %w", err)
		}

		log.Debug().Msgf("user: %+v", user)

		return nil
	})
	if err != nil {
		return model.MainPageData{}, fmt.Errorf("failed ReadCommitted: %w", err)
	}

	return model.MainPageData{Users: users, Rooms: rooms, User: user}, nil
}
