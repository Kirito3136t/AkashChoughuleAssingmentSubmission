package services

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockTransactionService struct {
	Queries          *database.Queries
	PortfolioService *PortfolioService
	StockService     *StockService
}

func NewStockTransactionService(queries *database.Queries, portfolioService *PortfolioService, s *StockService) *StockTransactionService {
	return &StockTransactionService{
		Queries:          queries,
		PortfolioService: portfolioService,
		StockService:     s,
	}
}

func (s *StockTransactionService) RegisterNewTransaction(ctx *gin.Context, transaction *models.TransactionRequest) (database.StockTransaction, error) {
	portfolioReq := &models.RecordPortfolioRequest{
		UserId:   transaction.UserID,
		StockId:  transaction.StockId,
		Quantity: transaction.Quantity,
		Type:     transaction.Type,
	}
	if _, err := s.PortfolioService.UpdateUserPortfolio(ctx, portfolioReq); err != nil {
		logger.Log.Error("Unable to update the users portfolio")
		return database.StockTransaction{}, err
	}

	params := database.CreateStockTransactionParams{
		ID:        uuid.New(),
		StockID:   transaction.StockId,
		UserID:    transaction.UserID,
		Quantity:  transaction.Quantity,
		Price:     transaction.Price,
		Type:      transaction.Type,
		CreatedAt: time.Now().UTC(),
	}

	return s.Queries.CreateStockTransaction(ctx, params)
}

func (s *StockTransactionService) RewardStock(ctx *gin.Context, req *models.RewardRequest) error {
	stock, err := s.StockService.GetStockBySymbol(ctx, req.StockSymbol)
	if err != nil {
		logger.Log.Error("Unable to update the users portfolio")
		return err
	}

	quantityFloat, _ := strconv.ParseFloat(req.Quantity, 64)
	valuationFloat, _ := strconv.ParseFloat(stock.Valuation, 64)
	amount := quantityFloat * valuationFloat

	transactionObject := models.TransactionRequest{
		UserID:   req.UserID,
		StockId:  stock.ID,
		Quantity: req.Quantity,
		Type:     "buy",
		Price:    fmt.Sprintf("%.4f", amount),
	}

	_, err = s.RegisterNewTransaction(ctx, &transactionObject)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to register the new transaction",
		})
		return err
	}

	return nil
}
