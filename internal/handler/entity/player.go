package entity

import (
	"net/http"

	"github.com/c-wide/albion-registry-api/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// GetPlayerInfo godoc
//
//	@Summary		Get basic player information
//	@Description	Get basic player information by player id
//	@Tags			entity
//	@Produce		json
//	@Param			region	path		string	true	"Server Region"
//	@Param			id		path		string	true	"Player ID"
//	@Success		200		{array}		database.GetPlayerRow
//	@Failure		404		{object}	echo.HTTPError
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/entity/player/{region}/{id} [get]
func (h *Handler) GetPlayerInfo(c echo.Context) error {
	var params BaseParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	player, err := h.queries.GetPlayer(c.Request().Context(), database.GetPlayerParams{
		Region:   params.Region,
		PlayerID: params.ID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Player not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	return c.JSON(http.StatusOK, player)
}
