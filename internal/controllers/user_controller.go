package controllers

import (
	"fmt"
	"net/http"

	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/database"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/logger"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/models"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/services"
	"github.com/Kirito3136t/AkashChoughuleAssingmentSubmission/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserService        *services.UserService
	TransactionService *services.StockTransactionService
	StockService       *services.StockService
	PortfolioService   *services.PortfolioService
}

func NewUserController(u *services.UserService, t *services.StockTransactionService, s *services.StockService, p *services.PortfolioService) *UserController {
	return &UserController{
		UserService:        u,
		TransactionService: t,
		StockService:       s,
		PortfolioService:   p,
	}
}

func (u *UserController) RegisterNewUser(ctx *gin.Context) {
	var req models.RegisterUserRequestObject
	stockSymbol := "TCS"
	quantity := "2.234"
	referralStockSymbol := "INFOSYS"
	referralQuantity := "3.21"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	_, err := u.UserService.GetUserByMail(ctx, req.Email)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User already exists",
		})
		return
	}

	if req.IsReferral {
		referralUser, err := u.UserService.GetUserByMail(ctx, req.ReferralUserEmail)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Unable to find the referral user",
			})
			return
		}

		rewardReq := models.RewardRequest{
			UserID:      referralUser.ID,
			StockSymbol: referralStockSymbol,
			Quantity:    referralQuantity,
		}

		err = u.TransactionService.RewardStock(ctx, &rewardReq)
		if err != nil {
			logger.Log.Error("Error: ", err)
			ctx.JSON(http.StatusBadGateway, gin.H{
				"error": "Unable to award the user a percent share of stock",
			})
			return
		}
	}

	user := database.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	newUser, err := u.UserService.RegisterUser(ctx, user)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to create a new user",
		})
		return
	}

	rewardReq := models.RewardRequest{
		UserID:      newUser.ID,
		StockSymbol: stockSymbol,
		Quantity:    quantity,
	}

	err = u.TransactionService.RewardStock(ctx, &rewardReq)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to award the user a percent share of stock",
		})
		return
	}

	logger.Log.Infof("New User created with ID: %v", newUser.ID)
	ctx.JSON(http.StatusAccepted, gin.H{
		"message": fmt.Sprintf("Congratulation you have been awarded with %v units of %v stocks", quantity, stockSymbol),
		"email":   newUser.Email,
		"name":    newUser.Name,
		"ID":      newUser.ID,
	})
}

func (u *UserController) LoginUser(ctx *gin.Context) {
	var req models.LoginUserRequestObject

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	user, err := u.UserService.GetUserByMail(ctx, req.Email)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No such user exist",
		})
		return
	}

	if user.Password != req.Password {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Unauthorized Access",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.Email)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to generate a token",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"token": token,
	})
}

func (u *UserController) UserPortfolio(ctx *gin.Context) {
	userId, err := services.ParseUserId(ctx)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to fetch the user id from token",
		})
		return
	}

	userPortfolio, err := u.PortfolioService.GetPortfolioByUserId(ctx, userId)
	if err != nil {
		logger.Log.Error("Error: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to fetch the users portfolio",
		})
		return
	}

	logger.Log.Info("Successfully retireved users portfolio")
	ctx.JSON(http.StatusOK, gin.H{
		"portfolio": userPortfolio,
	})
}
