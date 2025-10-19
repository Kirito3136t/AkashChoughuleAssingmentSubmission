package services

import (
	"time"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
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

// records a trasaction by an user
func (s *StockTransactionService) RegisterTransaction(ctx *gin.Context, transaction *models.TransactionRequestObject) (database.StockTransaction, error) {
	params := database.CreateStockTransactionParams{
		ID:              uuid.New(),
		StockID:         transaction.StockId,
		UserID:          transaction.UserID,
		Quantity:        transaction.Quantity,
		Price:           transaction.Price,
		Type:            transaction.Type,
		TransactionType: transaction.TransactionType,
		CreatedAt:       time.Now(),
	}

	return s.Queries.CreateStockTransaction(ctx, params)
}

func (s *StockTransactionService) GetUserTransactionForToday(ctx *gin.Context, user_id uuid.UUID, today_date time.Time) ([]database.StockTransaction, error) {
	transactionParams := database.GetTodaysUserStockParams{
		UserID:    user_id,
		CreatedAt: today_date,
	}

	return s.Queries.GetTodaysUserStock(ctx, transactionParams)
}
