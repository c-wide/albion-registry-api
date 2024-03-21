package handler

import (
	"net/http"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	zl "github.com/rs/zerolog/log"
)

type PlayersSearchParams struct {
	Name   string `query:"name" validate:"required,min=3"`
	Region string `query:"region" validate:"required,oneof=americas asia"`
}

func (h *Handler) PlayersSearch(c echo.Context) error {
	var params PlayersSearchParams
	if err := c.Bind(&params); err != nil {
		zl.Error().Err(err).Msg("Unable to bind players search params")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		zl.Error().Err(err).Msg("Invalid players search params")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	players, err := h.DB.FindPlayersByNameAndRegion(c.Request().Context(), database.FindPlayersByNameAndRegionParams{
		Name:   pgtype.Text{String: params.Name, Valid: true},
		Region: database.RegionEnum(params.Region),
	})
	if err != nil {
		zl.Error().Err(err).Msg("Unable to find players by name and region")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to find players by name and region")
	}

	return c.JSON(http.StatusOK, players)
}
