package history

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

// GetAllianceGuildHistory godoc
//
//	@Summary		Alliance guild history
//	@Description	Retrieve all guilds that have been part of the specified alliance
//	@Tags			history
//	@Produce		json
//	@Param			region		path		string	true	"Server Region"
//	@Param			alliance_id	path		string	true	"Alliance ID"
//	@Param			limit		query		int		false	"Limit (Default 10)"
//	@Param			offset		query		int		false	"Offset"
//	@Param			before_first_seen	query		string	false	"Cursor timestamp (RFC3339). Requires before_id"
//	@Param			before_id	query		string	false	"Cursor guild ID tiebreaker. Requires before_first_seen"
//	@Success		200			{array}		database.GetAllianceGuildHistoryRow
//	@Failure		400			{object}	echo.HTTPError
//	@Failure		500			{object}	echo.HTTPError
//	@Router			/history/alliance/{region}/{alliance_id}/guilds [get]
func (h *Handler) AllianceGuilds(c echo.Context) error {
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

	allianceHistory, err := h.queries.GetAllianceGuildHistory(c.Request().Context(), database.GetAllianceGuildHistoryParams{
		AllianceID: params.ID,
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

	return c.JSON(http.StatusOK, allianceHistory)
}
