package history

import (
	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/rs/zerolog"
)

type Handler struct {
	logger  zerolog.Logger
	queries *database.Queries
}

func New(logger zerolog.Logger, queries *database.Queries) *Handler {
	return &Handler{
		logger:  logger,
		queries: queries,
	}
}
