package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/syrlramadhan/api-bendahara-inovdes/config"
	"github.com/syrlramadhan/api-bendahara-inovdes/routes"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	appPort := os.Getenv("APP_PORT")
	fmt.Println("api running on http://localhost:" + appPort)

	db, err := config.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	routes.Routes(db, appPort)
}