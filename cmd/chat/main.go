package main

import (
	"context"
	"fmt"
	"github.com/patyukin/go-chat/internal/auth"
	"github.com/patyukin/go-chat/internal/cacher"
	"github.com/patyukin/go-chat/internal/config"
	"github.com/patyukin/go-chat/internal/db"
	"github.com/patyukin/go-chat/internal/dbconn"
	"github.com/patyukin/go-chat/internal/handler"
	"github.com/patyukin/go-chat/internal/server"
	"github.com/patyukin/go-chat/internal/server/router"
	"github.com/patyukin/go-chat/internal/usecase"
	"github.com/patyukin/go-chat/pkg/migrator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("unable to load config: %v", err)
	}

	dbConn, err := dbconn.New(ctx, cfg)
	if err != nil {
		log.Fatal().Msgf("failed connecting to db: %v", err)
	}

	err = migrator.UpMigrations(ctx, dbConn)
	if err != nil {
		log.Fatal().Msgf("failed migrating db: %v", err)
	}

	dbClient := db.New(dbConn)
	authClient := auth.New(cfg)

	log.Info().Msg("connected to db")

	chr, err := cacher.New(ctx, cfg)
	if err != nil {
		log.Fatal().Msgf("failed connecting to redis: %v", err)
	}

	uc := usecase.New(cfg, dbClient, authClient, chr)
	h := handler.New(uc)
	rtr := router.Init(h)

	errCh := make(chan error)

	srv := server.New(rtr)

	go func() {
		log.Info().Msgf("starting server on port: %d", cfg.HttpPort)
		if err = srv.Run(fmt.Sprintf(":%d", cfg.HttpPort)); err != nil {
			log.Error().Msgf("failed running http server: %v", err)
			errCh <- err
		}
	}()

	log.Info().Msg("Chat App Started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err = <-errCh:
		log.Error().Msgf("Failed to run, err: %v", err)
	case res := <-sigChan:
		if res == syscall.SIGINT || res == syscall.SIGTERM {
			log.Info().Msgf("Signal received")
		} else if res == syscall.SIGHUP {
			log.Info().Msgf("Signal received")
		}
	}

	log.Info().Msg("Shutting Down")

	if err = srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("failed server shutting down: %s", err.Error())
	}

	if err = dbClient.Close(); err != nil {
		log.Error().Msgf("failed db connection close: %s", err.Error())
	}
}
