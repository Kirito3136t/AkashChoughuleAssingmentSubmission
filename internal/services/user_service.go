package services

import (
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	Queries *database.Queries
}

func NewUserService(queries *database.Queries) *UserService {
	return &UserService{
		Queries: queries,
	}
}

func (u *UserService) GetUserByMail(ctx *gin.Context, email string) (database.User, error) {
	user, err := u.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return database.User{}, err
	}

	return user, nil
}

func (u *UserService) RegisterUser(ctx *gin.Context, user database.User) (database.User, error) {
	params := database.CreateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	return u.Queries.CreateUser(ctx, params)
}
