package stats

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Summary(c echo.Context) error {
	counts, err := h.queries.GetCountsOfEntities(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Stats summary unavailable")
	}

	return c.JSON(http.StatusOK, counts)
}
