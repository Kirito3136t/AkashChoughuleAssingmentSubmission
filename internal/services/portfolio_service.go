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

func (p *PortfolioService) FetcUserStock(ctx *gin.Context, user_id uuid.UUID, stock_id uuid.UUID) (database.Portfolio, error) {
	params := database.GetUserStockParams{
		UserID:  user_id,
		StockID: stock_id,
	}

	portfolio, err := p.Queries.GetUserStock(ctx, params)
	if err != nil {
		return database.Portfolio{}, err
	}

	return portfolio, nil
}

func (p *PortfolioService) UpdateUserPortfolio(ctx *gin.Context, portfolio *models.RecordPortfolioRequest) (database.Portfolio, error) {
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
	portfolio, err := p.Queries.GetPortfolioByUserId(ctx, userId)
	if err != nil {
		return []database.Portfolio{}, err
	}

	return portfolio, nil
}
