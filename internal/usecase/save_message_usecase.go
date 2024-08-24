package usecase

import (
	"context"
	"fmt"
)

func (u *UseCase) SaveMassageUseCase(ctx context.Context, roomUUID, userUUID, message string) error {
	err := u.registry.GetRepo().InsertMessage(ctx, roomUUID, userUUID, message)
	if err != nil {
		return fmt.Errorf("failed SaveMassageUseCase: %w", err)
	}

	return nil
}
