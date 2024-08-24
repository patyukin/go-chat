package cacher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/patyukin/go-chat/internal/config"
	"github.com/redis/go-redis/v9"
	"net"
)

type Cacher struct {
	redis *redis.Client
}

func New(ctx context.Context, cfg *config.Config) (*Cacher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Redis.Host, fmt.Sprintf("%d", cfg.Redis.Port)),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %w", err)
	}

	return &Cacher{redis: client}, nil
}

func (c *Cacher) SetConnection(ctx context.Context, connKey string, conn *websocket.Conn) error {
	connString, err := json.Marshal(conn)
	if err != nil {
		return fmt.Errorf("unable to marshal connection: %w", err)
	}

	_, err = c.redis.Set(ctx, connKey, connString, 0).Result()
	return err
}

func (c *Cacher) GetConnectionKeys(ctx context.Context, roomUUID string) ([]string, error) {
	pattern := "room:" + roomUUID + ":user:*"
	keys, err := c.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to get keys: %w", err)
	}

	return keys, nil
}

func (c *Cacher) DelConnection(ctx context.Context, connKey string) error {
	_, err := c.redis.Del(ctx, connKey).Result()
	return err
}

func (c *Cacher) GetConnection(ctx context.Context, connKey string) (*websocket.Conn, error) {
	connString, err := c.redis.Get(ctx, connKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to get connection: %w", err)
	}

	var conn *websocket.Conn
	err = json.Unmarshal([]byte(connString), &conn)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal connection: %w", err)
	}

	return conn, nil
}
