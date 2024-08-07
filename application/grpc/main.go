package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/sethvargo/go-envconfig"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/iainvm/deposits/application/grpc/gen/deposits/v1/depositsv1connect"
	"github.com/iainvm/deposits/application/grpc/handlers"
	"github.com/iainvm/deposits/common/postgres"
	"github.com/iainvm/deposits/internal/investors"
	store "github.com/iainvm/deposits/internal/investors/postgres"
)

type DBConfig struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT, default=5432"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"` // TODO: this would be pulled from a secrets vault
	Name     string `env:"NAME"`
}

type Config struct {
	Port     string   `env:"PORT, default=8080"`
	DBConfig DBConfig `env:", prefix=DB_"`
}

func main() {
	// Logger
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Parse Env Vars
	var config Config
	err := envconfig.Process(ctx, &config)
	if err != nil {
		logger.With("error", err).Error("failed to parse server configuration")
		panic(fmt.Errorf("failed to parse server configuration: %w", err))
	}
	logger.Debug("Config Processed", "Config", config)

	// DB
	dataSource := postgres.NewDataSource(
		config.DBConfig.Host,
		config.DBConfig.Port,
		config.DBConfig.User,
		config.DBConfig.Password,
		config.DBConfig.Name,
		false,
	)
	db, err := postgres.Connect(dataSource)
	if err != nil {
		logger.With("error", err).Error("failed to connect to DB")
		panic(fmt.Errorf("failed to connect to DB: %w", err))
	}
	logger.With("host", config.DBConfig.Host).With("port", config.DBConfig.Port).Info("Connected to DB")

	// Investors Handler
	investorsStore := store.NewStore(db)
	investorsService := investors.NewService(investorsStore)
	investorsServer := handlers.NewInvestorsHandler(investorsService, logger)

	// Register handlers
	mux := http.NewServeMux()
	path, handler := depositsv1connect.NewInvestorsServiceHandler(investorsServer)
	mux.Handle(path, handler)

	// Listen
	logger.With("port", config.Port).Info("Starting listener")
	http.ListenAndServe(
		fmt.Sprintf(":%s", config.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
