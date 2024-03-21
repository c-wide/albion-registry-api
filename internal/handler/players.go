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

	player, err := h.DB.GetPlayerByIdAndRegion(c.Request().Context(), database.GetPlayerByIdAndRegionParams{
		PlayerID: params.PlayerID,
		Region:   database.RegionEnum(params.Region),
	})
	if err != nil {
		zl.Error().Err(err).Msg("Unable to get player by ID")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get player by ID")
	}

	return c.JSON(http.StatusOK, player)
}

type PlayersHistoryParams struct {
	PlayerID string `query:"id" validate:"required"`
	Region   string `query:"region" validate:"required,oneof=americas asia"`
}

type PlayersHistoryResponse struct {
	PlayerGuildMemberships    []database.GetPlayerGuildMembershipsRow    `json:"guilds"`
	PlayerAllianceMemberships []database.GetPlayerAllianceMembershipsRow `json:"alliances"`
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

	playerGuildMemberships, err := h.DB.GetPlayerGuildMemberships(c.Request().Context(), database.GetPlayerGuildMembershipsParams{
		PlayerID: params.PlayerID,
		Region:   database.RegionEnum(params.Region),
	})
	if err != nil {
		zl.Error().Err(err).Msg("Unable to get player guild memberships")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get player guild memberships")
	}

	playerAllianceMemberships, err := h.DB.GetPlayerAllianceMemberships(c.Request().Context(), database.GetPlayerAllianceMembershipsParams{
		PlayerID: params.PlayerID,
		Region:   database.RegionEnum(params.Region),
	})
	if err != nil {
		zl.Error().Err(err).Msg("Unable to get player alliance memberships")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get player alliance memberships")
	}

	return c.JSON(http.StatusOK, PlayersHistoryResponse{
		PlayerGuildMemberships:    playerGuildMemberships,
		PlayerAllianceMemberships: playerAllianceMemberships,
	})
}
