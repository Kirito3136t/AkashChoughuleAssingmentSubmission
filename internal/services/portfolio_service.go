package services

import (
	"errors"
	"time"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PortfolioService struct {
	Queries *database.Queries
}

func NewPortfolioService(queries *database.Queries) *PortfolioService {
	return &PortfolioService{
		Queries: queries,
	}
}

func (p *PortfolioService) FetchUserPortfolioByStockId(ctx *gin.Context, user_id uuid.UUID, stock_id uuid.UUID) (database.Portfolio, error) {
	params := database.GetUserStockParams{
		UserID:  user_id,
		StockID: stock_id,
	}

	return p.Queries.GetUserStock(ctx, params)
}

func (p *PortfolioService) UpdateUserPortfolio(ctx *gin.Context, portfolio *models.RecordPortfolioRequest) (database.Portfolio, error) {
	logger.Log.Info(portfolio.Type)

	switch portfolio.Type {
	case "sell":
		portfolio.Quantity = "-" + portfolio.Quantity
	case "buy":
		//quantity stays positive
	default:
		logger.Log.Error("Invalid portfolio type. Please recheck")
		return database.Portfolio{}, errors.New("invalid portfolio type")
	}

	params := database.RecordUserPortfolioParams{
		ID:            uuid.New(),
		UserID:        portfolio.UserId,
		StockID:       portfolio.StockId,
		TotalQuantity: portfolio.Quantity,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return p.Queries.RecordUserPortfolio(ctx, params)
}

func (p *PortfolioService) GetPortfolioByUserId(ctx *gin.Context, userId uuid.UUID) ([]database.Portfolio, error) {
	return p.Queries.GetPortfolioByUserId(ctx, userId)
}
