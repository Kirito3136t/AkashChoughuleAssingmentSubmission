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
	stocks, err := s.Queries.GetAllStocks(ctx)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

// fetches stock by id
func (s *StockService) GetStockById(ctx *gin.Context, id uuid.UUID) (database.Stock, error) {
	stock, err := s.Queries.GetStockById(ctx, id)
	if err != nil {
		return database.Stock{}, err
	}

	return stock, nil
}

func (s *StockService) GetStockBySymbol(ctx *gin.Context, symbol string) (database.Stock, error) {
	stock, err := s.Queries.GetStockBySymbol(ctx, symbol)
	if err != nil {
		return database.Stock{}, err
	}

	return stock, nil
}
