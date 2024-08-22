package dbconn

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/patyukin/go-chat/internal/config"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	log.Info().Msg("connecting to db")
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Name,
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	log.Info().Msg("connected to db")

	err = dbConn.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	log.Info().Msg("pinged postgresql db")

	return dbConn, nil
}
