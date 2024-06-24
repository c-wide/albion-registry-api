package server

import (
	"context"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func New(logger zerolog.Logger, queries *database.Queries) *Server {
	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	registerMiddleware(e, logger)
	registerRoutes(e, logger, queries)

	return &Server{
		echo: e,
	}
}

func (s *Server) Start(address string) error {
	return s.echo.Start(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
