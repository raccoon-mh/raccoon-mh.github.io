package logger

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func logger(message string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s\n", currentTime, message)
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		uri := c.Request().RequestURI
		method := c.Request().Method
		resStatusCode := c.Response().Status
		logger(fmt.Sprintf("%s - %s - %s - %d", ip, method, uri, resStatusCode))
		return next(c)
	}
}
