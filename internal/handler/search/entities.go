package search

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

type SearchEntityParams struct {
	SearchTerm string `param:"search_term" validate:"required"`
	Region     string `param:"region" validate:"required,oneof=americas asia europe"`
	Limit      int32  `query:"limit" validate:"omitempty,min=1,max=50"`
}

// SearchEntities godoc
//
//	@Summary		Player, guild, or alliance lookup
//	@Description	Search for players, guilds, or alliances by their name or tag
//	@Tags			search
//	@Produce		json
//	@Param			region		path		string	true	"Server Region"
//	@Param			search_term	path		string	true	"Name or Tag"
//	@Param			limit		query		int		false	"Limit (Default 10)"
//	@Success		200			{array}		database.SearchEntitiesRow
//	@Failure		400			{object}	echo.HTTPError
//	@Failure		500			{object}	echo.HTTPError
//	@Router			/search/entities/{region}/{search_term} [get]
func (h *Handler) SearchEntities(c echo.Context) error {
	var params SearchEntityParams
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

	searchResults, err := h.queries.SearchEntities(c.Request().Context(), database.SearchEntitiesParams{
		Searchterm: params.SearchTerm,
		Region:     params.Region,
		Limit:      limit,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, searchResults)
}
