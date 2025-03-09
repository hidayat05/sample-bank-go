package app

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"sample-bank/app/startup"
	"sample-bank/config"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file. %v", err)
	}

	dbConfig := config.GetDBConfig()

	app := &startup.App{}
	grpcServer := app.Initialize(dbConfig)
	app.Run(grpcServer, ":"+os.Getenv("PORT"))
}
