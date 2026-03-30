package history

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

type PlayerGuildAllianceHistoryParams struct {
	BaseParams
	GuildID string `param:"guild" validate:"required"`
}

// GetPlayerGuildAllianceHistory godoc
//
//	@Summary		Player guild alliance history
//	@Description	Retrieve alliances the specified guild was a member of during the specified players tenure in that guild
//	@Tags			history
//	@Produce		json
//	@Param			region		path		string	true	"Server Region"
//	@Param			player_id	path		string	true	"Player ID"
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			limit		query		int		false	"Limit (Default 10)"
//	@Param			offset		query		int		false	"Offset"
//	@Param			before_first_seen	query		string	false	"Cursor timestamp (RFC3339). Requires before_id"
//	@Param			before_id	query		string	false	"Cursor alliance ID tiebreaker. Requires before_first_seen"
//	@Success		200			{array}		database.GetPlayerGuildAlliancesRow
//	@Failure		400			{object}	echo.HTTPError
//	@Failure		500			{object}	echo.HTTPError
//	@Router			/history/player/{region}/{player_id}/{guild_id}/alliances [get]
func (h *Handler) PlayerGuildAllianceHistory(c echo.Context) error {
	var params PlayerGuildAllianceHistoryParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cursor, err := parseCursorParams(params.BeforeFirstSeen, params.BeforeID, params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	playerHistory, err := h.queries.GetPlayerGuildAlliances(c.Request().Context(), database.GetPlayerGuildAlliancesParams{
		PlayerID:   params.ID,
		GuildID:    params.GuildID,
		Region:     params.Region,
		Limit:      defaultLimit(params.Limit),
		Offset:     cursor.Offset,
		UseCursor:  cursor.UseCursor,
		BeforeID:   cursor.BeforeID,
		BeforeTime: cursor.BeforeTime,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, playerHistory)
}
