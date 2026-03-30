package history

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

type PlayerHistoryParams struct {
	BaseParams
	AllianceLimit int32 `query:"allianceLimit" validate:"omitempty,min=1,max=20"`
}

// GetPlayerHistory godoc
//
//	@Summary		Player guild & alliance history
//	@Description	Retrieve all guilds the specified player has been a member of and the alliances the guild was a member of during the player's tenure
//	@Tags			history
//	@Produce		json
//	@Param			region			path		string	true	"Server Region"
//	@Param			player_id		path		string	true	"Player ID"
//	@Param			limit			query		int		false	"Limit (Default 10)"
//	@Param			offset			query		int		false	"Offset"
//	@Param			before_first_seen	query		string	false	"Cursor timestamp (RFC3339). Requires before_id"
//	@Param			before_id		query		string	false	"Cursor guild ID tiebreaker. Requires before_first_seen"
//	@Param			allianceLimit	query		int		false	"Alliance Limit (Default 5)"
//	@Success		200				{array}		database.GetPlayerHistoryRow
//	@Failure		400				{object}	echo.HTTPError
//	@Failure		500				{object}	echo.HTTPError
//	@Router			/history/player/{region}/{player_id} [get]
func (h *Handler) PlayerHistory(c echo.Context) error {
	var params PlayerHistoryParams
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

	limit := defaultLimit(params.Limit)

	allianceLimit := params.AllianceLimit
	if allianceLimit == 0 {
		allianceLimit = 5
	}

	playerHistory, err := h.queries.GetPlayerHistory(c.Request().Context(), database.GetPlayerHistoryParams{
		PlayerID:      params.ID,
		Region:        params.Region,
		Limit:         limit,
		Offset:        cursor.Offset,
		UseCursor:     cursor.UseCursor,
		BeforeID:      cursor.BeforeID,
		BeforeTime:    cursor.BeforeTime,
		AllianceLimit: allianceLimit,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, playerHistory)
}
