package main

import (
	"api/src/app"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	address string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	address = os.Getenv("ADDR") + ":" + os.Getenv("PORT")
}

func main() {
	printStartMessage()
	printAddress(address)
	app := app.App()
	if err := app.Start(address); err != nil {
		app.Logger.Fatal(err)
	}
}

func printStartMessage() {
	logo := `
                                       _           _ _   _        _      _     
 _ _ __ _ __ __ ___  ___ _ _ ___ _ __ | |_    __ _(_) |_| |_ _  _| |__  (_)___ 
| '_/ _` + "`" + ` / _/ _/ _ \/ _ \ ' \___| '  \| ' \ _/ _` + "`" + ` | |  _| ' \ || | '_ \_| / _ \
|_| \__,_\__\__\___/\___/_||_|  |_|_|_|_||_(_)__, |_|\__|_||_\_,_|_.__(_)_\___/
                                             |___/                             
================================================================================`
	fmt.Println(logo)
}

func printAddress(address string) {
	fmt.Println("SERVER strated at :", address)
}
