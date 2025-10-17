package services

import (
	"fmt"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/models"
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

func ParseUserId(ctx *gin.Context) (uuid.UUID, error) {
	userIdValue, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in token")
	}

	userIdStr, ok := userIdValue.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID in token is not a string")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	return userId, nil
}
