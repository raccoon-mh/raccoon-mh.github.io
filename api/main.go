package main

import (
	"api/src/app"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	address string
)

func init() {
	err := godotenv.Load("setup.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	address = os.Getenv("ADDR") + ":" + os.Getenv("PORT")
}

func main() {
	app := app.App()
	if err := app.Start(address); err != nil {
		app.Logger.Fatal(err)
	}
}
