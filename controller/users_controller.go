package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Queries *database.Queries
}

func NewUserController(q *database.Queries) *UserController {
	return &UserController{Queries: q}
}

func (u *UserController) registerUser(c *gin.Context) {
	context := c.Request.Context()

	var request models.RegisterUserRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = u.Queries.GetUserByEmail(context, request.Email)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("User already exists with email: %s", request.Email),
		})
		return
	}

}
