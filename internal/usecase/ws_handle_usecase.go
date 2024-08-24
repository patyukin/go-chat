package usecase

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/patyukin/go-chat/internal/handler/dto"
	"github.com/rs/zerolog/log"
	"maps"
)

func (u *UseCase) WsHandleUseCase(ctx context.Context, conn *websocket.Conn, roomUUID, userUUID string) error {
	user, err := u.registry.GetRepo().SelectUserByUUID(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("unable to get user: %w", err)
	}

	for {
		var message dto.SentMessage
		if err = conn.ReadJSON(&message); err != nil {
			log.Error().Msgf("unable to read message: %v", err)
			break
		}

		log.Debug().Msgf("message: %+v, roomID id: %s", message, roomUUID)
		if err = u.SaveMassageUseCase(ctx, roomUUID, userUUID, message.Message); err != nil {
			log.Error().Msgf("unable to save message: %v", err)
			break
		}

		log.Debug().Msgf("message saved")

		var keys []string
		keys, err = u.chr.GetConnectionKeys(ctx, roomUUID)
		if err != nil {
			return fmt.Errorf("unable to get connection keys: %w", err)
		}

		log.Debug().Msgf("keys: %+v", keys)

		broadcastMessage := map[string]string{
			"sender":  user.Login,
			"content": message.Message,
		}

		for _, key := range keys {
			log.Debug().Msgf("key: %s", key)
			var currentConnection *websocket.Conn
			currentConnection, err = u.chr.GetConnection(ctx, key)
			if err != nil {
				return fmt.Errorf("unable to get connection: %w", err)
			}

			log.Debug().Msgf("currentConnection: %+v", currentConnection)
			connErr := currentConnection.WriteJSON(maps.Clone(broadcastMessage))
			if connErr != nil {
				log.Error().Msgf("unable to broadcast message: %v", connErr)
				if connCLoseErr := conn.Close(); connCLoseErr != nil {
					return fmt.Errorf("unable to close connection: %w", connCLoseErr)
				}

				if delErr := u.chr.DelConnection(ctx, key); delErr != nil {
					return fmt.Errorf("unable to delete connection: %w", delErr)
				}
			}

			log.Debug().Msgf("message broadcasted")
		}
	}

	return nil
}
