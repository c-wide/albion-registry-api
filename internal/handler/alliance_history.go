package handler

import (
	"net/http"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

func (h *Handler) AllianceHistory(c echo.Context) error {
	var params HistoryParams
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

	allianceHistory, err := h.db.GetAllianceHistory(c.Request().Context(), database.GetAllianceHistoryParams{
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
