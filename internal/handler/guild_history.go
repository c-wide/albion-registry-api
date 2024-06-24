package handler

import (
	"net/http"

	"github.com/ao-tools/albion-registry-api/internal/database"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GuildHistory(c echo.Context) error {
	var params HistoryParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	guildHistory, err := h.db.GetGuildHistory(c.Request().Context(), database.GetGuildHistoryParams{
		GuildID:       params.ID,
		Region:        params.Region,
		Alliancelimit: 10,
		Playerlimit:   20,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "An error occurred while processing your request")
	}

	formattedHistory := GuildHistoryResponse{
		Alliances: make([]database.GetGuildAllianceHistoryRow, 0),
		Players:   make([]database.GetGuildPlayerHistoryRow, 0),
	}

	for _, row := range guildHistory {
		switch row.Type {
		case "alliance":
			formattedHistory.Alliances = append(formattedHistory.Alliances, database.GetGuildAllianceHistoryRow{
				AllianceID: row.ID,
				Name:       row.Name,
				Tag:        row.Tag,
				FirstSeen:  row.FirstSeen,
				LastSeen:   row.LastSeen,
			})
		case "player":
			formattedHistory.Players = append(formattedHistory.Players, database.GetGuildPlayerHistoryRow{
				PlayerID:  row.ID,
				Name:      *row.Name,
				FirstSeen: row.FirstSeen,
				LastSeen:  row.LastSeen,
			})
		}
	}

	return c.JSON(http.StatusOK, formattedHistory)
}
