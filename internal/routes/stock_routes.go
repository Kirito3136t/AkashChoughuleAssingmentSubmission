package routes

import (
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/controllers"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/middleware"
	"github.com/gin-gonic/gin"
)

func StockRoutes(r *gin.Engine, stockController *controllers.StockController) {
	stockGroup := r.Group("/stocks")
	stockGroup.Use(middleware.AuthMiddleware())
	{
		stockGroup.GET("/", stockController.GetAllStocks)
	}
}
