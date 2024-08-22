package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type QueryExecutor interface {
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...interface{}) *sql.Row
}

type Client struct {
	db *sql.DB
}

func (c *Client) GetRepo() *Repository {
	return &Repository{
		db: c.db,
	}
}

type Handler func(ctx context.Context, repo *Repository) error

func New(db *sql.DB) *Client {
	return &Client{db: db}
}

func (c *Client) ReadCommitted(ctx context.Context, f Handler) error {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			if errRollback := tx.Rollback(); errRollback != nil {
				slog.ErrorContext(ctx, "failed to rollback transaction", slog.Any("error", errRollback))
			}
		}
	}()

	repo := &Repository{db: tx}

	if err = f(ctx, repo); err != nil {
		return fmt.Errorf("failed to execute handler: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (r *Repository) InsertIntoUsers(ctx context.Context, id, login string) error {
	timeNow := time.Now().UTC()
	query := `INSERT INTO users (auth_user_id, login, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, id, login, timeNow)
	if err != nil {
		return fmt.Errorf("failed insert into users: %w", err)
	}

	return nil
}
