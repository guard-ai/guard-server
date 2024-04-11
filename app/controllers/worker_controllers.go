package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/guard-ai/guard-server/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (controller *Controller) Record(e echo.Context) error {
	type Request struct {
		Logs   []models.Log   `json:"logs"`
		Events []models.Event `json:"events"`
	}
	request := Request{}
	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	ctx := context.Background()
	conn, err := controller.db.Acquire(ctx)
	if err != nil {
		e.Logger().Error(err)
		return e.NoContent(http.StatusInternalServerError)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		e.Logger().Error(err)
		return e.NoContent(http.StatusInternalServerError)
	}

	for _, log := range request.Logs {
		_, err = tx.Exec(ctx, `INSERT INTO public."Logs" (id, region, utterance) VALUES ($1, $2, $3)`, log.Id, log.Region, log.Utterance)
		if err != nil {
			e.Logger().Error(err)
			if err := tx.Rollback(ctx); err != nil {
				e.Logger().Error(err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("there was an error trying to create a log with id: %v", log.Id))
		}
	}

	for _, event := range request.Events {
		_, err = tx.Exec(ctx, `INSERT INTO Public."Events" (id, level, location, category, log_id, description) VALUES ($1, $2, ST_GeomFromGeoJSON($3), $4, $5, $6)`, event.Id, event.Level, event.Location.AsGeoJSON(), event.Category, event.LogId, event.Description)
		if err != nil {
			e.Logger().Error(err)
			if err := tx.Rollback(ctx); err != nil {
				e.Logger().Error(err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("there was an error trying to create an event with id: %v", event.Id))
		}

		if err := tx.Commit(ctx); err != nil {
			e.Logger().Error(err)
			return e.NoContent(http.StatusInternalServerError)
		}

		go func(event models.Event) {
			ctx := context.Background()
			conn, err := controller.db.Acquire(ctx)
			if err != nil {
				e.Logger().Error(err)
				return
			}

			rows, err := conn.Query(ctx, `
				SELECT push_token 
				FROM Public."Users"
				WHERE ST_DWithin(location::geography, ST_GeomFromGeoJSON($1)::geography, 5 * 1609.34)`, event.Location.AsGeoJSON())

			users := []string{}

			for rows.Next() {
				var pushToken string
				if err := rows.Scan(&pushToken); err != nil {
					e.Logger().Error(err)
					continue
				}
				users = append(users, pushToken)
			}

			if err := controller.notifier.Broadcast(event, users); err != nil {
				fmt.Println(err.Error())
			}
		}(event)
	}

	return e.NoContent(http.StatusOK)
}
