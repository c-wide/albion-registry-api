package handler

import (
	"errors"
	"net/http"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	zl "github.com/rs/zerolog/log"
)

type PlayersDetailsParams struct {
	PlayerID string `query:"id" validate:"required"`
	Region   string `query:"region" validate:"required,oneof=americas asia"`
}

func (h *Handler) PlayersDetails(c echo.Context) error {
	var params PlayersDetailsParams
	if err := c.Bind(&params); err != nil {
		zl.Error().Err(err).Msg("Unable to bind players details params")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		zl.Error().Err(err).Msg("Invalid players details params")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	player, err := h.DB.GetPlayerDetails(c.Request().Context(), database.GetPlayerDetailsParams{
		PlayerID: params.PlayerID,
		Region:   database.RegionEnum(params.Region),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Player not found")
		}

		zl.Error().Err(err).Msg("Unable to get player by ID")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get player by ID")
	}

	return c.JSON(http.StatusOK, player)
}
