package entity

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// GetGuildInfo godoc
//
//	@Summary		Get basic guild information
//	@Description	Get basic guild information by guild id
//	@Tags			entity
//	@Produce		json
//	@Param			region	path		string	true	"Server Region"
//	@Param			id		path		string	true	"Guild ID"
//	@Success		200		{array}		database.GetGuildRow
//	@Failure		404		{object}	echo.HTTPError
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/entity/guild/{region}/{id} [get]
func (h *Handler) GetGuildInfo(c echo.Context) error {
	var params BaseParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	guild, err := h.queries.GetGuild(c.Request().Context(), database.GetGuildParams{
		Region:  params.Region,
		GuildID: params.ID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Guild not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, guild)
}
