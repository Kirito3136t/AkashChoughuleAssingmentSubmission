package routes

import (
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/controllers"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userController *controllers.UserController) {
	userGroup := r.Group("/api/users")
	{
		userGroup.POST("/register", userController.RegisterNewUser)
		userGroup.POST("/login", userController.LoginUser)
	}
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.GET("/portfolio", userController.UserPortfolio)
	}
}
