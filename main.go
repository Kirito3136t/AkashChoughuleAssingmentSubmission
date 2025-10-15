package main

import (
	"os"
	"strconv"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/db"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/logger"
	routes "github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() (port, dbName, dbUser, dbPass, dbHost string, dbPort int) {
	port = os.Getenv("PORT")
	if port == "" {
		logger.Log.Warn("PORT not set, using default 8080")
		port = "8080"
	}

	dbName = os.Getenv("DB_NAME")
	if dbName == "" {
		logger.Log.Fatal("DB_NAME not set")
	}

	dbUser = os.Getenv("DB_USERNAME")
	if dbUser == "" {
		logger.Log.Fatal("DB_USERNAME not set")
	}

	dbPass = os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		logger.Log.Fatal("DB_PASSWORD not set")
	}

	dbHost = os.Getenv("DB_HOST")
	if dbHost == "" {
		logger.Log.Fatal("DB_HOST not set")
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		logger.Log.Warn("DB_PORT not set, using default 5432")
		portStr = "5432"
	}

	var err error
	dbPort, err = strconv.Atoi(portStr)
	if err != nil {
		logger.Log.Fatal("Invalid DB_PORT value: ", err)
	}

	return
}

func main() {
	logger.InitLogger()

	if err := godotenv.Load(".env"); err != nil {
		logger.Log.Warn("Unable to load .env file: ", err)
	}

	port, dbName, dbUser, dbPass, dbHost, dbPort := loadEnv()

	conn := db.Connect(dbName, dbUser, dbPass, dbHost, dbPort)

	queries := database.New(conn)

	router := gin.New()
	routes.StockRoutes(router, queries)

	logger.Log.Info("Server starting at port: ", port)
	router.Run(":" + port)
}
