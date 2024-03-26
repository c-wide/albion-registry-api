package handler

import (
	"net/http"
	"time"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
	zl "github.com/rs/zerolog/log"
)

type PlayersHistoryParams struct {
	PlayerID string `query:"id" validate:"required"`
	Region   string `query:"region" validate:"required,oneof=americas asia"`
}

type PlayersHistoryGuild struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

type PlayersHistoryAlliance struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tag       string    `json:"tag"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

type PlayersHistoryResponse struct {
	Guilds    []PlayersHistoryGuild    `json:"guilds"`
	Alliances []PlayersHistoryAlliance `json:"alliances"`
}

const recordLimit = 100

func (h *Handler) PlayersHistory(c echo.Context) error {
	var params PlayersHistoryParams
	if err := c.Bind(&params); err != nil {
		zl.Error().Err(err).Msg("Unable to bind players history params")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		zl.Error().Err(err).Msg("Invalid players history params")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	playerExists, err := h.DB.DoesPlayerExist(c.Request().Context(), database.DoesPlayerExistParams{
		PlayerID: params.PlayerID,
		Region:   database.RegionEnum(params.Region),
	})
	if err != nil {
		zl.Error().Err(err).Msg("Unable to check if player exists")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to check if player exists")
	}

	if !playerExists {
		return echo.NewHTTPError(http.StatusNotFound, "Player not found")
	}

	playerHistory, err := h.DB.GetPlayerHistory(c.Request().Context(), database.GetPlayerHistoryParams{
		PlayerID: params.PlayerID,
		Region:   database.RegionEnum(params.Region),
		Limit:    recordLimit,
	})
	if err != nil {
		zl.Error().Err(err).Msg("Unable to get player history")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get player history")
	}

	guilds := make([]PlayersHistoryGuild, 0)
	alliances := make([]PlayersHistoryAlliance, 0)

	for _, row := range playerHistory {
		switch row.Type {
		case "guild":
			guilds = append(guilds, PlayersHistoryGuild{
				ID:        row.ID,
				Name:      row.Name,
				FirstSeen: row.FirstSeen,
				LastSeen:  row.LastSeen,
			})
		case "alliance":
			alliances = append(alliances, PlayersHistoryAlliance{
				ID:        row.ID,
				Name:      row.Name,
				Tag:       row.Tag.(string),
				FirstSeen: row.FirstSeen,
				LastSeen:  row.LastSeen,
			})
		}
	}

	return c.JSON(http.StatusOK, PlayersHistoryResponse{
		Guilds:    guilds,
		Alliances: alliances,
	})
}
