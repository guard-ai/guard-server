package server

import (
	"net/http"

	"github.com/guard-ai/guard-server/app/controllers"
	"github.com/guard-ai/guard-server/pkg"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, controller *controllers.Controller) {
	publicRoutes(e, controller)
	privateRoutes(e, controller)
}

func publicRoutes(e *echo.Echo, controller *controllers.Controller) {
	e.GET("/events/near/:uuid", controller.EventsNear)

	e.POST("/user/ping", controller.PingUser)
	e.POST("/user", controller.CreateUser)

	e.GET("*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
}

func privateRoutes(e *echo.Echo, controller *controllers.Controller) {
	e.POST("/worker/record", controller.Record, pkg.AuthenticatedAsWorkerMiddleware())
}
