package pkg

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthenticatedAsWorkerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token != fmt.Sprintf("Bearer %v", Env().WorkerAuthToken) {
				c.Error(fmt.Errorf("unauthorized to make this request\n"))
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}
