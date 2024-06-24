package handler

import "github.com/ao-tools/albion-registry-api/internal/database"

type Handler struct {
	db *database.Queries
}

type HistoryParams struct {
	ID     string `query:"id" validate:"required"`
	Region string `query:"region" validate:"required,oneof=americas asia europe"`
	Limit  int32  `query:"limit" validate:"omitempty,min=1,max=50"`
	Offset int32  `query:"offset" validate:"omitempty,min=0"`
}

type GuildHistoryResponse struct {
	Alliances []database.GetGuildAllianceHistoryRow `json:"alliances"`
	Players   []database.GetGuildPlayerHistoryRow   `json:"players"`
}
