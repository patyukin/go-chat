package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/patyukin/go-chat/internal/model"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	db QueryExecutor
}

func (r *Repository) SelectUsersWithoutUserUUID(ctx context.Context, id string) ([]model.User, error) {
	query := `SELECT u.id, u.login, u.auth_user_id FROM users AS u WHERE id != $1 order by u.login`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed select all users: %w", err)
	}

	var users []model.User
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.Login, &user.AuthUserID); err != nil {
			return nil, fmt.Errorf("failed scan row: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) SelectAllUsers(ctx context.Context) ([]model.User, error) {
	query := `SELECT id, login, auth_user_id FROM users order by login`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed select all users: %w", err)
	}

	var users []model.User
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.Login, &user.AuthUserID); err != nil {
			return nil, fmt.Errorf("failed scan row: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) SelectAllRooms(ctx context.Context, id string) ([]model.Room, error) {
	query := `
SELECT 
    r.id,
    r.name
FROM rooms AS r 
    JOIN users_rooms AS ur ON r.id = ur.room_id
WHERE ur.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed select all users: %w", err)
	}

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		if err = rows.Scan(&room.ID, &room.Name); err != nil {
			return nil, fmt.Errorf("failed scan room row: %w", err)
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *Repository) SelectUserUUIDByAuthUserID(ctx context.Context, id string) (string, error) {
	query := `SELECT id FROM users WHERE auth_user_id = $1`
	var uuid string
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&uuid); err != nil {
		return "", fmt.Errorf("failed select login: %w", err)
	}

	return uuid, nil
}

func (r *Repository) InsertIntoRooms(ctx context.Context, name string) (string, error) {
	query := `INSERT INTO rooms (name) VALUES ($1) RETURNING id`
	var id string
	if err := r.db.QueryRowContext(ctx, query, name).Scan(&id); err != nil {
		return "", fmt.Errorf("failed insert into rooms: %w", err)
	}

	return id, nil
}

func (r *Repository) InsertIntoRoomsUsers(ctx context.Context, roomID string, users []string) error {
	query := `INSERT INTO users_rooms (user_id, room_id) VALUES ($1, $2)`

	for _, user := range users {
		if _, err := r.db.ExecContext(ctx, query, user, roomID); err != nil {
			return fmt.Errorf("failed insert into users_rooms: %w", err)
		}
	}

	return nil
}

func (r *Repository) SelectMessagesByRoomUUIDUserUUID(ctx context.Context, userUUID, roomUUID string) ([]model.Message, error) {
	query := `
SELECT 
	m.id,
	m.room_id,
	m.user_id,
	m.content,
	m.created_at,
	u.login
FROM messages AS m
	JOIN users_rooms AS ur ON m.room_id = ur.room_id
	JOIN users AS u ON m.user_id = u.id
WHERE ur.user_id = $1 AND m.room_id = $2
ORDER BY m.created_at
`

	rows, err := r.db.QueryContext(ctx, query, userUUID, roomUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed select messages: %w", err)
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error().Msgf("failed close rows: %v", err)
		}
	}(rows)

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed rows: %w", err)
	}

	var messages []model.Message
	for rows.Next() {
		var message model.Message
		if err = rows.Scan(&message.ID, &message.RoomUUID, &message.UserUUID, &message.Content, &message.CreatedAt, &message.User.Login); err != nil {
			return nil, fmt.Errorf("failed scan message row: %w", err)
		}

		messages = append(messages, message)
	}

	log.Debug().Msgf("messages: %+v", messages)

	return messages, nil
}

func (r *Repository) SelectRoomByRoomUUID(ctx context.Context, roomUUID string) (model.Room, error) {
	query := `SELECT id, name, created_at FROM rooms WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, roomUUID)
	var room model.Room
	if err := row.Scan(&room.ID, &room.Name, &room.CreatedAt); err != nil {
		return model.Room{}, fmt.Errorf("failed select room: %w", err)
	}

	return room, nil
}

func (r *Repository) InsertMessage(ctx context.Context, roomUUID, userUUID, message string) error {
	query := `INSERT INTO messages (room_id, user_id, content) VALUES ($1, $2, $3)`

	if _, err := r.db.ExecContext(ctx, query, roomUUID, userUUID, message); err != nil {
		return fmt.Errorf("failed insert message: %w", err)
	}

	return nil
}

func (r *Repository) SelectUserByUUID(ctx context.Context, userUUID string) (model.User, error) {
	query := `SELECT id, login FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userUUID)

	var user model.User
	if err := row.Scan(&user.ID, &user.Login); err != nil {
		return model.User{}, fmt.Errorf("failed select user: %w", err)
	}

	return user, nil
}
