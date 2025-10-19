package services

import (
	"time"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService struct {
	Queries *database.Queries
}

func NewUserService(queries *database.Queries) *UserService {
	return &UserService{
		Queries: queries,
	}
}

// fetches user by email
func (u *UserService) GetUserByMail(ctx *gin.Context, email string) (database.User, error) {
	return u.Queries.GetUserByEmail(ctx, email)
}

// registers a new user
func (u *UserService) RegisterUser(ctx *gin.Context, user *models.RequestBodyRegisterUser) (database.User, error) {
	params := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return u.Queries.CreateUser(ctx, params)
}
