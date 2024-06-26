package server

import (
	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/c-wide/albion-registry-api/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/swaggo/echo-swagger"

	_ "github.com/c-wide/albion-registry-api/third_party/swagger"
)

func registerRoutes(e *echo.Echo, logger zerolog.Logger, queries *database.Queries) {
	// Initialize handler
	h := handler.New(logger, queries)

	// Stats routes
	statsGroup := e.Group("/stats")
	statsGroup.GET("/summary", h.Stats.Summary)

	// History routes
	historyGroup := e.Group("/history")
	historyGroup.GET("/player/:region/:id", h.History.PlayerHistory)
	historyGroup.GET("/player/:region/:id/:guild/alliances", h.History.PlayerGuildAllianceHistory)
	historyGroup.GET("/guild/:region/:id/alliances", h.History.GuildAlliances)
	historyGroup.GET("/guild/:region/:id/players", h.History.GuildPlayers)
	historyGroup.GET("/alliance/:region/:id/guilds", h.History.AllianceGuilds)

	// Search routes
	searchGroup := e.Group("/search")
	searchGroup.GET("/entities/:region", h.Search.SearchEntities)

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
