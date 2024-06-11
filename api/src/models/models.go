package models

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db *gorm.DB
)

func init() {
	var err error
	dbPath := os.Getenv("DBPATH")
	Db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println(err.Error())
		panic("fail to connect DB")
	}
}

func DbMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tx := Db.WithContext(c.Request().Context())
		c.Set("tx", tx)
		return next(c)
	}
}
