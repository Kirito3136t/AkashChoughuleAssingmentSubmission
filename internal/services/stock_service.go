package services

import (
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockService struct {
	Queries *database.Queries
}

func NewStockService(queries *database.Queries) *StockService {
	return &StockService{
		Queries: queries,
	}
}

// fetches all the stocks
func (s *StockService) GetAllStocks(ctx *gin.Context) ([]database.Stock, error) {
	return s.Queries.GetAllStocks(ctx)
}

// fetches stock by id
func (s *StockService) GetStockById(ctx *gin.Context, id uuid.UUID) (database.Stock, error) {
	return s.Queries.GetStockById(ctx, id)
}

// fetches stock by symbol
func (s *StockService) GetStockBySymbol(ctx *gin.Context, symbol string) (database.Stock, error) {
	return s.Queries.GetStockBySymbol(ctx, symbol)
}
