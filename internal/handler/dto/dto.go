package dto

import "github.com/patyukin/go-chat/internal/model"

type CreateRoomV1Request struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
}

type CreateRoomV1Response struct {
	RoomID string `json:"id"`
}

type PageRoomV1Response struct {
	Messages []model.Message `json:"messages"`
	Users    []model.User    `json:"users"`
	Room     model.Room      `json:"room"`
	User     model.User      `json:"user"`
}

type SentMessage struct {
	Message string `json:"content"`
}
