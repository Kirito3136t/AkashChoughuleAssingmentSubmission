package controllers

import (
	"net/http"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/services"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/utils"
	"github.com/gin-gonic/gin"
)

type StockController struct {
	StockService       *services.StockService
	TransactionService *services.StockTransactionService
	PortfolioService   *services.PortfolioService
}

func NewStockController(stockService *services.StockService, transactionService *services.StockTransactionService, portfolioService *services.PortfolioService) *StockController {
	return &StockController{
		StockService:       stockService,
		TransactionService: transactionService,
		PortfolioService:   portfolioService,
	}
}

func (s *StockController) GetAllStocks(ctx *gin.Context) {
	stocks, err := s.StockService.GetAllStocks(ctx)
	if err != nil {
		logger.Log.Error("FUNC(GetAllStocks): Erro fetching all stocks :", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch the stocks",
		})
		return
	}

	logger.Log.Info("FUNC(GetAllStocks): Retrieved all the stocks")
	ctx.JSON(http.StatusOK, gin.H{
		"data": utils.MapDatabaseStocks(stocks),
	})
}
