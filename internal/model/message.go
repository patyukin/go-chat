package model

type Message struct {
	ID        string `json:"id,omitempty"`
	Content   string `json:"content"`
	UserUUID  string `json:"user_id"`
	RoomUUID  string `json:"room_id"`
	CreatedAt string `json:"created_at"`
	Room      Room
	User      User
}
