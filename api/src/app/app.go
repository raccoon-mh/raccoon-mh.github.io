package app

import (
	"api/src/models"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	app     *echo.Echo
	appOnce sync.Once
)

func App() *echo.Echo {
	appOnce.Do(func() {
		app = echo.New()
		app.Use(models.DbMiddleware)

		app.GET("/alive", func(c echo.Context) error {
			return c.JSON(http.StatusOK, models.CommonResponseStatusOK("alive"))
		})

	})
	return app
}
