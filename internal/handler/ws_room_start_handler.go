package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/patyukin/go-chat/internal/handler/dto"
	"github.com/rs/zerolog/log"
	"maps"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WsRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomUUID := r.URL.Path[len("/ws/rooms/"):]
	userUUID := r.Header.Get("X-User-UUID")
	user, err := h.uc.SelectUserByUUID(r.Context(), userUUID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Msgf("unable to get user uuid")
		http.Error(w, "Unable to render page", http.StatusBadRequest)
		return
	}

	log.Debug().Msgf("roomUUID: %s, userUUID: %s", roomUUID, userUUID)

	connKey := "room:" + roomUUID + ":user:" + userUUID
	h.connections.Store(connKey, conn)

	defer func() {
		h.connections.Delete(connKey)
	}()

	for {
		var message dto.SentMessage
		if err = conn.ReadJSON(&message); err != nil {
			log.Error().Msgf("unable to read message: %v", err)
			break
		}

		log.Debug().Msgf("message: %+v, roomID id: %s", message, roomUUID)
		if err = h.uc.SaveMassageUseCase(r.Context(), roomUUID, userUUID, message.Message); err != nil {
			log.Error().Msgf("unable to save message: %v", err)
			break
		}

		log.Debug().Msgf("message saved")

		broadcastMessage := map[string]string{
			"sender":  user.Login,
			"content": message.Message,
		}

		h.connections.Range(func(key, value interface{}) bool {
			currentConnectionKey := key.(string)
			if con, ok := value.(*websocket.Conn); ok && isUserRoom(currentConnectionKey, roomUUID) {
				err = con.WriteJSON(maps.Clone(broadcastMessage))
				if err != nil {
					log.Error().Msgf("unable to broadcast message: %v", err)
					if connCLoseErr := con.Close(); connCLoseErr != nil {
						log.Error().Msgf("unable to close connection: %v", connCLoseErr)
					}
				}
			}

			return true
		})

		log.Debug().Msgf("message broadcasted")
	}
}

func isUserRoom(connKey, roomUUID string) bool {
	return fmt.Sprintf("room:%s:", roomUUID) == connKey[:len(fmt.Sprintf("room:%s:", roomUUID))]

}
