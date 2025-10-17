package main

import (
	"os"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/app"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Info("No .env file located")
	}

	a := app.NewApp()
	port := os.Getenv("PORT")
	if port == "" {
		logger.Log.Info("Unable to locate PORT variable")
	}

	a.Run(":" + port)
}
