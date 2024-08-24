package usecase

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) ValidateTokenUseCase(ctx context.Context, token string) (string, error) {
	log.Info().Msgf("token: %v", token)
	result, err := u.authClient.ValidateToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("failed ValidateTokenUseCase: %w", err)
	}

	log.Info().Msgf("user uuid: %v", result)

	uuid, err := u.registry.GetRepo().SelectUserUUIDByAuthUserID(ctx, result.UUID)
	if err != nil {
		return "", fmt.Errorf("failed SelectLoginByAuthUserID: %w", err)
	}

	return uuid, nil
}
