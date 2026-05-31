package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ktravels-publicsite-api/graph"
	"ktravels-publicsite-api/internal/features"
	"ktravels-publicsite-api/internal/platform/config"
	"ktravels-publicsite-api/internal/platform/database"
	gqlserver "ktravels-publicsite-api/internal/platform/graphql"
	"ktravels-publicsite-api/internal/platform/logger"
	rmq "ktravels-publicsite-api/internal/platform/rabbitmq"

	currencyfeature "ktravels-publicsite-api/internal/features/currency"
	curapp "ktravels-publicsite-api/internal/features/currency/application"
	exchangeratefeature "ktravels-publicsite-api/internal/features/exchange_rate"
	erapp "ktravels-publicsite-api/internal/features/exchange_rate/application"
	erinfra "ktravels-publicsite-api/internal/features/exchange_rate/infrastructure"
	providersvc "ktravels-publicsite-api/internal/features/provider_service"
	sharedcurrency "ktravels-publicsite-api/internal/shared/currency"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.App.Name, cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	db, err := database.New(ctx, cfg.Database.URI, cfg.Database.Database)
	if err != nil {
		log.Error().Err(err).Str("operation", "Database.Connect").Msg("failed to connect to database")
		os.Exit(1)
	}
	defer db.Close(ctx)

	rabbitClient, err := rmq.New(rmq.Config{
		HostName: cfg.RabbitMQ.HostName,
		UserName: cfg.RabbitMQ.UserName,
		Password: cfg.RabbitMQ.Password,
		Port:     cfg.RabbitMQ.Port,
	})
	if err != nil {
		log.Error().Err(err).Str("operation", "RabbitMQ.Connect").Msg("failed to connect to RabbitMQ")
		os.Exit(1)
	}
	defer rabbitClient.Close()

	rateCache := erapp.NewRateCache()
	currencyCache := curapp.NewCurrencyCache()
	_ = sharedcurrency.NewCurrencyConverter(rateCache, currencyCache)

	erRepo := erinfra.NewMongoRepository(db.Database())

	deps := features.Deps{
		RabbitMQ: rabbitClient,
		Log:      log,
	}

	featureList := []features.Feature{
		providersvc.NewFeature(db.Database()),
		exchangeratefeature.NewFeature(db.Database(), rateCache),
		currencyfeature.NewFeature(db.Database(), currencyCache),
	}

	for _, f := range featureList {
		if err := f.Register(ctx, deps); err != nil {
			log.Error().Err(err).Str("operation", "Feature.Register").Msg("failed to register feature")
			os.Exit(1)
		}
	}

	log.Info().Str("operation", "App.Start").Str("environment", cfg.App.Env).Msg("starting application")

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	gqlHandler := gqlserver.NewHandler(
		graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			ExchangeRateRepo: erRepo,
		}}),
		cfg.GraphQL.PlaygroundEnabled,
	)
	server := &http.Server{
		Addr:         addr,
		Handler:      gqlHandler,
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
		log.Error().Err(err).Str("operation", "HTTP.ListenAndServe").Msg("server error")
	case sig := <-sigChan:
		log.Info().Str("operation", "App.Signal").Str("signal", sig.String()).Msg("received signal")
	}

	log.Info().Str("operation", "HTTP.Shutdown").Msg("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Str("operation", "HTTP.Shutdown").Msg("server shutdown error")
	}
}
