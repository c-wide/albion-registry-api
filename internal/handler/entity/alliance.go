package entity

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// GetAllianceInfo godoc
//
//	@Summary		Get basic alliance information
//	@Description	Get basic alliance information by alliance id
//	@Tags			entity
//	@Produce		json
//	@Param			region	path		string	true	"Server Region"
//	@Param			id		path		string	true	"Alliance ID"
//	@Success		200		{array}		database.GetAllianceRow
//	@Failure		404		{object}	echo.HTTPError
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/entity/alliance/{region}/{id} [get]
func (h *Handler) GetAllianceInfo(c echo.Context) error {
	var params BaseParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	alliance, err := h.queries.GetAlliance(c.Request().Context(), database.GetAllianceParams{
		Region:     params.Region,
		AllianceID: params.ID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Alliance not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, alliance)
}
