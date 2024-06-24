package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) StatsSummary(c echo.Context) error {
	stats, err := h.db.GetCountsOfEntities(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get stats summary")
	}

	return c.JSON(http.StatusOK, stats)
}
