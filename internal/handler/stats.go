package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	zl "github.com/rs/zerolog/log"
)

func (h *Handler) StatsSummary(c echo.Context) error {
	stats, err := h.DB.GetCountsOfEntities(c.Request().Context())
	if err != nil {
		zl.Error().Err(err).Msg("Unable to get stats summary")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get stats summary")
	}

	return c.JSON(http.StatusOK, stats)
}
