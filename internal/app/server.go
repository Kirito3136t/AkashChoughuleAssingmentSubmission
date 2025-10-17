package app

import (
	"os"
	"strconv"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/db"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/controllers"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/routes"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/services"
	"github.com/gin-gonic/gin"
)

type App struct {
	Router  *gin.Engine
	Queries *database.Queries
}

func loadEnv() (dbName, dbUser, dbPass, dbHost string, dbPort int) {
	dbName = os.Getenv("DB_NAME")
	if dbName == "" {
		logger.Log.Error("DB_NAME not set")
	}

	dbUser = os.Getenv("DB_USERNAME")
	if dbUser == "" {
		logger.Log.Error("DB_USERNAME not set")
	}

	dbPass = os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		logger.Log.Error("DB_PASSWORD not set")
	}

	dbHost = os.Getenv("DB_HOST")
	if dbHost == "" {
		logger.Log.Error("DB_HOST not set")
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		logger.Log.Warn("DB_PORT not set, using default 5432")
		portStr = "5432"
	}

	var err error
	dbPort, err = strconv.Atoi(portStr)
	if err != nil {
		logger.Log.Error("Invalid DB_PORT value: ", err)
	}

	return
}

func NewApp() *App {
	dbName, dbUser, dbPass, dbHost, dbPort := loadEnv()

	router := gin.New()

	router.Use(gin.Recovery())
	dbConn := db.Connect(dbName, dbUser, dbPass, dbHost, dbPort)

	dbQueries := database.New(dbConn)

	//services
	stockService := services.NewStockService(dbQueries)
	portfolioService := services.NewPortfolioService(dbQueries)
	userService := services.NewUserService(dbQueries)
	transactionService := services.NewStockTransactionService(dbQueries, portfolioService, stockService)

	//controllers
	stockController := controllers.NewStockController(stockService, transactionService, portfolioService)
	userController := controllers.NewUserController(userService, transactionService, stockService, portfolioService)

	//routes
	routes.StockRoutes(router, stockController)
	routes.UserRoutes(router, userController)

	app := &App{
		Router:  router,
		Queries: dbQueries,
	}

	return app
}

func (a *App) Run(addr string) error {
	logger.Log.Infof("Starting server on %s", addr)
	return a.Router.Run(addr)
}
