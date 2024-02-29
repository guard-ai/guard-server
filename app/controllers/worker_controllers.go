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
		fmt.Printf("%v\n", log)
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
		fmt.Printf("%v\n", event)
		_, err = tx.Exec(ctx, `INSERT INTO Public."Events" (id, level, location, category, log_id) VALUES ($1, $2, ST_GeomFromGeoJSON($3), $4, $5)`, event.Id, event.Level, event.Location.AsGeoJSON(), event.Category, event.LogId)
		if err != nil {
			e.Logger().Error(err)
			if err := tx.Rollback(ctx); err != nil {
				e.Logger().Error(err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("there was an error trying to create an event with id: %v", event.Id))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		e.Logger().Error(err)
		return e.NoContent(http.StatusInternalServerError)
	}

	// TODO: send users notifications

	return e.NoContent(http.StatusOK)
}
