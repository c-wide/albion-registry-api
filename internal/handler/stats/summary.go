package stats

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// StatsSummary godoc
//
//	@Summary		General API statistics
//	@Description	Retrieves total number of tracked players, guilds, and alliances
//	@Tags			stats
//	@Produce		json
//	@Success		200	{array}		database.GetCountsOfEntitiesRow
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/stats/summary [get]
func (h *Handler) Summary(c echo.Context) error {
	counts, err := h.queries.GetCountsOfEntities(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Stats summary unavailable")
	}

	return c.JSON(http.StatusOK, counts)
}
