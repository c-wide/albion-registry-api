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

	limit := params.Limit
	if limit == 0 {
		limit = 10
	}

	allianceHistory, err := h.queries.GetAllianceGuildHistory(c.Request().Context(), database.GetAllianceGuildHistoryParams{
		AllianceID: params.ID,
		Region:     params.Region,
		Limit:      limit,
		Offset:     params.Offset,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, allianceHistory)
}
