package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/models"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch the stocks",
		})
		return
	}

	logger.Log.Info("Fetched all the stocks")
	ctx.JSON(http.StatusOK, gin.H{
		"stocks": services.MapDatabaseStocks(stocks),
	})
}

func (s *StockController) UserActionOnStock(ctx *gin.Context) {
	var req models.RequestBodyTransaction

	userId, err := services.ParseUserId(ctx)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to fetch the user id from token",
		})
		return
	}

	stockIdParam := ctx.Param("stock_id")
	stock_id, err := uuid.Parse(stockIdParam)
	if err != nil {
		logger.Log.Error("Unable to parse the stock id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please validate the api url",
		})
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Error("Unable to parse the request body", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please validate the request body",
		})
		return
	}

	if req.Type == "sell" {
		portfolio, err := s.PortfolioService.FetcUserStock(ctx, userId, stock_id)
		if err != nil {
			logger.Log.Error("Error fetching user's portfolio:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "You do not own this stock",
			})
			return
		}

		logger.Log.Info(portfolio)

		portfolio_quantity, err := strconv.ParseFloat(portfolio.TotalQuantity, 64)
		if err != nil {
			logger.Log.Error("Error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Unable to parse the float value",
			})
			return
		}

		if req.Quantity > portfolio_quantity {
			logger.Log.Errorf("You only have %v quantity of this stock", portfolio.TotalQuantity)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "You do not have enough stocks to sell",
			})
			return
		}

	}

	stock, err := s.StockService.GetStockById(ctx, stock_id)
	if err != nil {
		logger.Log.Error("Unable to sretrieve the stock details")
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("Cannot retrive details about stock : %v", stock_id),
		})
		return
	}

	stockValuation, err := strconv.ParseFloat(stock.Valuation, 64)
	if err != nil {
		logger.Log.Error("Unable to parse the string value to float")
		return
	}
	price := stockValuation * req.Quantity

	transactionObject := models.TransactionRequest{
		UserID:   userId,
		StockId:  stock_id,
		Quantity: fmt.Sprintf("%.6f", req.Quantity),
		Type:     req.Type,
		Price:    fmt.Sprintf("%.4f", price),
	}

	transaction, err := s.TransactionService.RegisterNewTransaction(ctx, &transactionObject)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to register the new transaction",
		})
		return
	}

	logger.Log.Info("Allocation of the stock successful")
	ctx.JSON(http.StatusAccepted, gin.H{
		"transaction": transaction,
	})
}
