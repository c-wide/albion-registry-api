package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/ao-tools/albion-registry-api/internal/handler"
	adapter "github.com/ao-tools/albion-registry-api/internal/zerolog"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
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
		zl.Info().Msg("Unable to load environment variables from .env file")
	}

	// Create a new Axiom writer
	writer, err := adapter.New()
	if err != nil {
		zl.Fatal().Err(err).Msg("Unable to create Axiom writer")
	}

	// Defer the close of the writer, very important
	defer writer.Close()

	// Create a new logger and set the default logger to use it
	zl.Logger = zerolog.New(io.MultiWriter(writer, os.Stderr)).With().Timestamp().Logger()

	// Database stuff
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		zl.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer conn.Close(ctx)

	db := database.New(conn)

	// Echo new echo instance
	e := echo.New()

	// Create and register validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Register Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	// Initialize handler
	h := &handler.Handler{DB: db}

	// Register groups
	stats := e.Group("/stats")
	players := e.Group("/players")

	// Register routes
	stats.GET("/summary", h.StatsSummary)
	players.GET("/search", h.PlayersSearch)
	players.GET("/details", h.PlayersDetails)
	players.GET("/history", h.PlayersHistory)

	// Start the server
	err = e.Start(fmt.Sprintf(":%s", os.Getenv("DEFAULT_PORT")))
	if err != nil {
		zl.Fatal().Err(err).Msg("Unable to start server")
	}
}
