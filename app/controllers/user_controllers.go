package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/guard-ai/guard-server/pkg/models"
	"github.com/labstack/echo/v4"
)

func (c *Controller) CreateUser(e echo.Context) error {
	type Request struct {
		Id        uuid.UUID    `json:"id"`
		Location  models.Point `json:"location"`
		PushToken string       `json:"push_token"`
	}
	request := Request{}
	if err := e.Bind(&request); err != nil {
		e.Logger().Error(err)
		return e.NoContent(http.StatusInternalServerError)
	}

	ctx := context.Background()
	conn, err := c.db.Acquire(ctx)
	if err != nil {
		e.Logger().Error(e)
		return e.NoContent(http.StatusInternalServerError)
	}

	fmt.Printf("Creating user at location: %s\n", request.Location.AsGeoJSON())
	_, err = conn.Exec(ctx, `INSERT INTO Public."Users" (id, location, push_token) VALUES ($1, ST_GeomFromGeoJSON($2), $3)`, request.Id, request.Location.AsGeoJSON(), request.PushToken)
	if err != nil {
		e.Logger().Error(err)
		return e.NoContent(http.StatusInternalServerError)
	}

	return e.NoContent(http.StatusOK)
}

func (c *Controller) PingUser(e echo.Context) error {
	type PingRequest struct {
		Id       uuid.UUID    `json:"id"`
		Location models.Point `json:"location"`
	}
	request := PingRequest{}
	if err := e.Bind(&request); err != nil {
		e.Logger().Error(err)
		return e.NoContent(http.StatusInternalServerError)
	}

	ctx := context.Background()
	conn, err := c.db.Acquire(ctx)
	if err != nil {
		e.Logger().Error(e)
		return e.NoContent(http.StatusInternalServerError)
	}

	_, err = conn.Exec(ctx, `UPDATE INTO Public."Users" SET id = $1, location = ST_GeomFromGeoJSON($2), last_ping = NOW() WHERE id = $1`, request.Id, request.Location.AsGeoJSON())
	if err != nil {
		e.Logger().Error(err)
		return e.NoContent(http.StatusInternalServerError)
	}

	return e.NoContent(http.StatusOK)
}
