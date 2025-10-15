package controllers

import (
	"net/http"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/gin-gonic/gin"
)

type StockController struct {
	Queries *database.Queries
}

func NewStockController(q *database.Queries) *StockController {
	return &StockController{Queries: q}
}

func (s *StockController) GetAllStocks(c *gin.Context) {
	context := c.Request.Context()

	stocks, err := s.Queries.GetAllStocks(context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetched the stocks",
		})

		return
	}

	if stocks == nil {
		stocks = []database.Stock{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stocks,
	})
}
