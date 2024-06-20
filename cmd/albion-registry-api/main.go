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
	zl "github.com/rs/zerolog/log"
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

	// TODO: do stuff from here

	// Echo new echo instance
	e := echo.New()

	// Create and register validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Register Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	// Initialize handler
	h := &handler.Handler{DB: queries}

	// Register groups
	stats := e.Group("/stats")
	players := e.Group("/players")

	// Register routes
	stats.GET("/summary", h.StatsSummary)
	players.GET("/history", h.PlayersHistory)

	// Start the server
	err = e.Start(fmt.Sprintf(":%s", os.Getenv("DEFAULT_PORT")))
	if err != nil {
		zl.Fatal().Err(err).Msg("Unable to start server")
	}
}
