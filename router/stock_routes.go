package routes

import (
	controllers "github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/controller"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/gin-gonic/gin"
)

func StockRoutes(router *gin.Engine, queries *database.Queries) {
	stockController := controllers.NewStockController(queries)

	router.GET("/stocks", stockController.GetAllStocks)
}
