package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
	zl "github.com/rs/zerolog/log"
)

type PlayersHistoryParams struct {
	PlayerID string `query:"id" validate:"required"`
	Region   string `query:"region" validate:"required,oneof=americas asia europe"`
}

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
		Region:   params.Region,
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
		Region:   params.Region,
	})
	if err != nil {
		zl.Error().Err(err).Msg("Unable to get player history")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get player history")
	}

	var history []interface{}
	err = json.Unmarshal(playerHistory, &history)
	if err != nil {
		zl.Error().Err(err).Msg("Unable to unmarshal player history")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to unmarshal player history")
	}

	return c.JSON(http.StatusOK, history)
}
