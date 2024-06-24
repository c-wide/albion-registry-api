package history

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

type PlayerGuildAllianceHistoryParams struct {
	BaseParams
	GuildID string `param:"guild" validate:"required"`
}

func (h *Handler) PlayerGuildAllianceHistory(c echo.Context) error {
	var params PlayerGuildAllianceHistoryParams
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

	playerHistory, err := h.queries.GetPlayerGuildAlliances(c.Request().Context(), database.GetPlayerGuildAlliancesParams{
		PlayerID: params.ID,
		GuildID:  params.GuildID,
		Region:   params.Region,
		Limit:    limit,
		Offset:   params.Offset,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, playerHistory)
}
