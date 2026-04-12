package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ktravels-publicsite-api/internal/platform/config"
	"ktravels-publicsite-api/internal/platform/database"
	"ktravels-publicsite-api/internal/platform/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	db, err := database.New(ctx, cfg.Database.URI, cfg.Database.Database)
	if err != nil {
		log.Error().Msgf("failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close(ctx)

	log.Info().Msgf("starting %s in %s mode", cfg.App.Name, cfg.App.Env)

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Error().Msgf("server error: %v", err)
	case sig := <-sigChan:
		log.Info().Msgf("received signal: %v", sig)
	}

	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msgf("server shutdown error: %v", err)
	}

	_ = db
}
