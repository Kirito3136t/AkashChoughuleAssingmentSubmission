package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/models"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MapDatabaseStocks(dbStocks []database.Stock) []models.StockResponseObject {
	stocks := make([]models.StockResponseObject, len(dbStocks))
	for i, s := range dbStocks {
		stocks[i] = models.StockResponseObject{
			ID:          s.ID.String(),
			StockSymbol: s.StockSymbol,
			Valuation:   s.Valuation,
		}
	}

	return stocks
}

func MapTransactions(ctx *gin.Context, requestedTransactionType string, s *services.StockService, transactions []database.StockTransaction) ([]models.ResponseBodyTransaction, error) {

	var responseTransactions []models.ResponseBodyTransaction

	for _, t := range transactions {
		stock, err := s.GetStockById(ctx, t.StockID)
		if err != nil {
			logger.Log.Error("func(MapTransactions): Unable to fetch the stock: ", err)
			return nil, fmt.Errorf("invalid stock id: %v", t.StockID)
		}

		if t.TransactionType == requestedTransactionType {
			responseTransactions = append(responseTransactions, models.ResponseBodyTransaction{
				UserID:          t.UserID.String(),
				StockSymbol:     stock.StockSymbol,
				Quantity:        t.Quantity,
				Price:           t.Price,
				Type:            t.Type,
				TransactionType: t.TransactionType,
			})
		}
	}

	return responseTransactions, nil
}

func MapPortfolio(ctx *gin.Context, s *services.StockService, portfolio []database.Portfolio) []models.ResponseObjectPortfolio {
	responsePortfolio := make([]models.ResponseObjectPortfolio, len(portfolio))
	for i, p := range portfolio {
		stock, err := s.GetStockById(ctx, p.StockID)
		if err != nil {
			logger.Log.Error("func(MapPortfolio): Unable to fetch the stock: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Stock id is not valid",
			})
			return []models.ResponseObjectPortfolio{}
		}

		quantity, _ := strconv.ParseFloat(p.TotalQuantity, 64)
		stockValuation, _ := strconv.ParseFloat(stock.Valuation, 64)

		currentValuation := quantity * stockValuation
		responsePortfolio[i] = models.ResponseObjectPortfolio{
			ID:               p.ID,
			StockSymbol:      stock.StockSymbol,
			Quantity:         p.TotalQuantity,
			CurrentValuation: currentValuation,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
		}
	}

	return responsePortfolio
}

func ParseUserId(ctx *gin.Context) (uuid.UUID, error) {
	userIdValue, exists := ctx.Get("user_id")

	if !exists {
		logger.Log.Errorf("func(ParseUserId): Couldn't fetch the user id from token")
		return uuid.Nil, fmt.Errorf("user id not found in token")
	}

	userIdStr, ok := userIdValue.(string)
	if !ok {
		logger.Log.Errorf("func(ParseUserId): Unable to parse the userId to string")
		return uuid.Nil, fmt.Errorf("user ID is not in string format")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		logger.Log.Errorf("func(ParseUserId): cannot convert token to uuid")
		return uuid.Nil, fmt.Errorf("invalid user ID format: %v", err)
	}

	return userId, nil
}

func RewardStock(ctx *gin.Context, stockService services.StockService, transactionService services.StockTransactionService, portfolioService services.PortfolioService, userId uuid.UUID, rewardType string) {
	var stockSymbol string
	var quantity string

	switch rewardType {
	case "registration":
		stockSymbol = "TCS"
		quantity = "2.32198"

	case "referral":
		stockSymbol = "ITC"
		quantity = "1.23532"
	case "milestone":
		stockSymbol = "INFOSYS"
		quantity = "0.7832"
	}

	stock, err := stockService.GetStockBySymbol(ctx, stockSymbol)
	if err != nil {
		logger.Log.Error("func(RewardStock): Unable to retrieve the stock", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to retrieve the stock",
		})
	}

	stockValuation, _ := strconv.ParseFloat(stock.Valuation, 64)
	quantityFloat, _ := strconv.ParseFloat(quantity, 64)

	currentValuation := stockValuation * quantityFloat

	object := models.TransactionRequestObject{
		UserID:          userId,
		StockId:         stock.ID,
		Type:            "buy",
		Quantity:        quantity,
		Price:           fmt.Sprintf("%.4f", currentValuation),
		TransactionType: "reward",
	}

	portfolioRequest := models.RecordPortfolioRequest{
		UserId:   userId,
		StockId:  stock.ID,
		Quantity: quantity,
		Type:     "buy",
	}

	_, err = portfolioService.UpdateUserPortfolio(ctx, &portfolioRequest)
	if err != nil {
		logger.Log.Error("func(UserActionOnStock): Cannot record the portfolio: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to update the users portfolio",
		})
		return
	}

	transactionService.RegisterTransaction(ctx, &object)
}
