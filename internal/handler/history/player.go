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

func (h *Handler) PlayerHistory(c echo.Context) error {
	var params PlayerHistoryParams
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

	allianceLimit := params.AllianceLimit
	if allianceLimit == 0 {
		allianceLimit = 5
	}

	playerHistory, err := h.queries.GetPlayerHistory(c.Request().Context(), database.GetPlayerHistoryParams{
		PlayerID:      params.ID,
		Region:        params.Region,
		Limit:         limit,
		Offset:        params.Offset,
		Alliancelimit: allianceLimit,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, playerHistory)
}
