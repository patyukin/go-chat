package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/go-chat/internal/model"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) SignUpUseCase(ctx context.Context, in model.SignUpRequest) error {
	result, err := u.authClient.SignUp(ctx, in)
	if err != nil {
		return fmt.Errorf("failed RegisterUser: %w", err)
	}

	log.Info().Msgf("user uuid: %v", result)

	if err = u.registry.GetRepo().InsertIntoUsers(ctx, result.UUID, in.Login); err != nil {
		return fmt.Errorf("failed save id to pg: %w", err)
	}

	return nil
}
