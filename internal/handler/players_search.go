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
	Limit  int32  `query:"limit" validate:"omitempty,min=1"`
	Offset int32  `query:"offset" validate:"omitempty,min=0"`
}

const defaultRecordLimit = 100

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

	if params.Limit == 0 {
		params.Limit = defaultRecordLimit
	}

	players, err := h.DB.SearchPlayers(c.Request().Context(), database.SearchPlayersParams{
		Name:   pgtype.Text{String: params.Name, Valid: true},
		Region: database.RegionEnum(params.Region),
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		zl.Error().Err(err).Msg("Unable to find players by name and region")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to find players by name and region")
	}

	return c.JSON(http.StatusOK, players)
}
