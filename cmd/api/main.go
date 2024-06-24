package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	adapter "github.com/axiomhq/axiom-go/adapters/zerolog"
	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/c-wide/albion-registry-api/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

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

	// Create server
	s := server.New(logger, queries)

	// Create SIGINT context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		err = s.Start(fmt.Sprintf(":%s", os.Getenv("DEFAULT_PORT")))
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Shutting down the server")
		}
	}()

	// Wait for SIGINT signal
	<-ctx.Done()

	// Gracefully shutdown server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Unable to shut down server gracefully")
	}

}
