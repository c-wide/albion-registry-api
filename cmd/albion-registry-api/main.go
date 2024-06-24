package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/ao-tools/albion-registry-api/internal/handler"
	adapter "github.com/axiomhq/axiom-go/adapters/zerolog"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Unable to load environment variables from .env file")
	}

	// Create default logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// If Axiom logging is enabled, create new logger
	_, axiomEnabled := os.LookupEnv("AXIOM_TOKEN")
	if axiomEnabled {
		writer, err := adapter.New()
		if err != nil {
			logger.Fatal().Err(err).Msg("Unable to create Axiom writer")
		}

		defer writer.Close()

		logger = zerolog.New(io.MultiWriter(writer, os.Stderr)).With().Timestamp().Logger()
	}

	// Database stuff
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to create database pool")
	}
	defer pool.Close()

	queries := database.New(pool)

	// Create new echo instance
	e := echo.New()

	// Create and register validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Register Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:      true,
		LogRemoteIP:     true,
		LogHost:         true,
		LogMethod:       true,
		LogURI:          true,
		LogUserAgent:    true,
		LogStatus:       true,
		LogError:        true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var event *zerolog.Event

			if v.Error != nil {
				event = logger.Error().Err(v.Error)
			} else {
				event = logger.Info()
			}

			event.
				Dur("latency", v.Latency).
				Str("remoteIP", v.RemoteIP).
				Str("host", v.Host).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("userAgent", v.UserAgent).
				Int("status", v.Status).
				Int64("responseSize", v.ResponseSize)

			event.Send()

			return nil
		},
	}))

	// Initialize handler
	h := handler.New(queries)

	// Register groups
	statsGroup := e.Group("/stats")
	historyGroup := e.Group("/history")

	// Register routes
	statsGroup.GET("/summary", h.StatsSummary)
	historyGroup.GET("/players", h.PlayerHistory)
	historyGroup.GET("/guilds", h.GuildHistory)
	historyGroup.GET("/alliances", h.AllianceHistory)

	// Start the server
	err = e.Start(fmt.Sprintf(":%s", os.Getenv("DEFAULT_PORT")))
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to start server")
	}
}
