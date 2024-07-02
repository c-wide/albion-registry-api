package handler

import (
	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/c-wide/albion-registry-api/internal/handler/entity"
	"github.com/c-wide/albion-registry-api/internal/handler/history"
	"github.com/c-wide/albion-registry-api/internal/handler/search"
	"github.com/c-wide/albion-registry-api/internal/handler/stats"
	"github.com/rs/zerolog"
)

type Handler struct {
	Stats   *stats.Handler
	History *history.Handler
	Search  *search.Handler
	Entity  *entity.Handler
}

func New(logger zerolog.Logger, queries *database.Queries) *Handler {
	return &Handler{
		Stats:   stats.New(logger, queries),
		History: history.New(logger, queries),
		Search:  search.New(logger, queries),
		Entity:  entity.New(logger, queries),
	}
}
