package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/go-chat/internal/model"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) SignInUseCase(ctx context.Context, in model.SignInRequest) (model.SignInResponse, error) {
	tokens, err := u.authClient.SignIn(ctx, in)
	if err != nil {
		return model.SignInResponse{}, fmt.Errorf("failed authClient.SignIn from SignIn: %w", err)
	}

	log.Debug().Msgf("tokens: %v", tokens)

	return tokens, nil
}
