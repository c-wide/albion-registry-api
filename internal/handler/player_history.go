package handler

import (
	"net/http"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

func (h *Handler) PlayerHistory(c echo.Context) error {
	var params HistoryParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	playerHistory, err := h.db.GetPlayerHistory(c.Request().Context(), database.GetPlayerHistoryParams{
		PlayerID: params.ID,
		Region:   params.Region,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, playerHistory)
}
