package model

type User struct {
	ID         string `json:"id"`
	Login      string `json:"login"`
	AuthUserID string `json:"auth_user_id"`
}
