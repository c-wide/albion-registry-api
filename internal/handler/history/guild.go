package history

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

// GetGuildAllianceHistory godoc
//
//	@Summary		Guild alliance history
//	@Description	Retrieve all alliances that the specified guild has been a member of
//	@Tags			history
//	@Produce		json
//	@Param			region		path		string	true	"Server Region"
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			limit		query		int		false	"Limit (Default 10)"
//	@Param			offset		query		int		false	"Offset"
//	@Param			before_first_seen	query		string	false	"Cursor timestamp (RFC3339). Requires before_id"
//	@Param			before_id	query		string	false	"Cursor alliance ID tiebreaker. Requires before_first_seen"
//	@Success		200			{array}		database.GetGuildAllianceHistoryRow
//	@Failure		400			{object}	echo.HTTPError
//	@Failure		500			{object}	echo.HTTPError
//	@Router			/history/guild/{region}/{guild_id}/alliances [get]
func (h *Handler) GuildAlliances(c echo.Context) error {
	var params BaseParams
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

	guildHistory, err := h.queries.GetGuildAllianceHistory(c.Request().Context(), database.GetGuildAllianceHistoryParams{
		GuildID:    params.ID,
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

	return c.JSON(http.StatusOK, guildHistory)
}

// GetGuildPlayerHistory godoc
//
//	@Summary		Guild player history
//	@Description	Retrieve all players that have been a member of the specified guild
//	@Tags			history
//	@Produce		json
//	@Param			region		path		string	true	"Server Region"
//	@Param			guild_id	path		string	true	"Guild ID"
//	@Param			limit		query		int		false	"Limit (Default 10)"
//	@Param			offset		query		int		false	"Offset"
//	@Param			before_first_seen	query		string	false	"Cursor timestamp (RFC3339). Requires before_id"
//	@Param			before_id	query		string	false	"Cursor player ID tiebreaker. Requires before_first_seen"
//	@Success		200			{array}		database.GetGuildPlayerHistoryRow
//	@Failure		400			{object}	echo.HTTPError
//	@Failure		500			{object}	echo.HTTPError
//	@Router			/history/guild/{region}/{guild_id}/players [get]
func (h *Handler) GuildPlayers(c echo.Context) error {
	var params BaseParams
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

	guildHistory, err := h.queries.GetGuildPlayerHistory(c.Request().Context(), database.GetGuildPlayerHistoryParams{
		GuildID:    params.ID,
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

	return c.JSON(http.StatusOK, guildHistory)
}
