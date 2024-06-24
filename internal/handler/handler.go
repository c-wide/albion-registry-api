package handler

import (
	"github.com/ao-tools/albion-registry-api/internal/database"
)

func New(queries *database.Queries) *Handler {
	return &Handler{
		db: queries,
	}
}
