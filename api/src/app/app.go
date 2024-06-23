package app

import (
	"api/src/handler/logger"
	"api/src/handler/models"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	app     *echo.Echo
	appOnce sync.Once
)

func App() *echo.Echo {
	appOnce.Do(func() {
		app = echo.New()
		app.HideBanner = true
		app.HidePort = true

		app.Use(models.DbMiddleware)
		app.Use(logger.LoggerMiddleware)
		app.Use(middleware.Recover())

		app.GET("/", root)
		app.GET("/alive", alive)

	})
	return app
}

func root(c echo.Context) error {
	return c.JSON(http.StatusOK, models.CommonResponseStatusOK("raccoon-mh's playground api server."))
}

func alive(c echo.Context) error {
	return c.JSON(http.StatusOK, models.CommonResponseStatusOK("alive"))
}
