package cacher

import (
	"context"
	"fmt"
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
