package routes

import (
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/controllers"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userController *controllers.UserController) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", userController.RegisterNewUser)
		userGroup.POST("/login", userController.LoginUser)
	}
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.POST("/transact/:stock_id", userController.UserActionOnStock)
		userGroup.GET("/today-stocks/:user_id", userController.FetchUsersTransactionForToday)
		userGroup.GET("/today-rewards/:user_id", userController.FetchUsersRewardForToday)
		userGroup.GET("/portfolio", userController.FetchUserPortfolio)
	}
}
