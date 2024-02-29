package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/guard-ai/guard-server/pkg/models"
	"github.com/labstack/echo/v4"
)

func (controller *Controller) EventsNear(e echo.Context) error {
	ctx := context.Background()
	conn, err := controller.db.Acquire(ctx)
	if err != nil {
		e.Logger().Error(err)
		log.Printf("%v\n", err)
		return e.NoContent(http.StatusInternalServerError)
	}
	defer conn.Release()

	id := e.Param("uuid")
	row := conn.QueryRow(ctx, `SELECT id, ST_AsGeoJSON(location) FROM Public."Users" WHERE id = $1`, id)
	user := models.User{}
	if err := row.Scan(&user.Id, &user.Location); err != nil {
		e.Logger().Error(e)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("there was an error trying to find user with id: %v", id))
	}

	rows, err := conn.Query(ctx, `
	SELECT id, level, ST_AsGeoJSON(location), category, log_id, created_at
	FROM Public."Events"
	WHERE ST_DWithin(location::geography, ST_GeomFromGeoJSON($1)::geography, 5 * 1609.34)`, user.Location.AsGeoJSON())
	if err != nil {
		e.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "there was an error trying to get all events")
	}
	defer rows.Close()

	events := []models.Event{}

	for rows.Next() {
		event := models.Event{}
		err := rows.Scan(&event.Id, &event.Level, &event.Location, &event.Category, &event.LogId, &event.CreatedAt)
		if err != nil {
			e.Logger().Error(err)
			continue
		}
		events = append(events, event)
	}

	return e.JSON(http.StatusOK, echo.Map{
		"events": events,
	})
}
