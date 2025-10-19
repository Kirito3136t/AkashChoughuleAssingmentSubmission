package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	var req *models.RequestBodyRegisterUser

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Func(RegisterNewUser): Invalid request body :", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	_, err := u.UserService.GetUserByMail(ctx, req.Email)
	if err == nil {
		logger.Log.Error("Func(RegisterNewUser): User already exists", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User already exists",
		})
		return
	}

	if req.IsReferral {
		user, err := u.UserService.GetUserByMail(ctx, req.ReferralUserEmail)
		if err != nil {
			logger.Log.Error("Func(RegisterNewUser): No referral user exists: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Unable to find the referral user",
			})
			return
		}

		utils.RewardStock(ctx, *u.StockService, *u.TransactionService, *u.PortfolioService, user.ID, "referral")
	}

	newUser, err := u.UserService.RegisterUser(ctx, req)
	if err != nil {
		logger.Log.Error("Func(RegisterNewUser): Error creating new user: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to create a new user",
		})
		return
	}

	utils.RewardStock(ctx, *u.StockService, *u.TransactionService, *u.PortfolioService, newUser.ID, "registration")

	logger.Log.Infof("Func(RegisterNewUser): New User created with ID: %v", newUser.ID)
	ctx.JSON(http.StatusAccepted, gin.H{
		"email": newUser.Email,
		"name":  newUser.Name,
		"ID":    newUser.ID,
	})
}

func (u *UserController) LoginUser(ctx *gin.Context) {
	var req models.RequestBodyLoginUser

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("func(LoginUser): Inavlid request body: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	user, err := u.UserService.GetUserByMail(ctx, req.Email)
	if err != nil {
		logger.Log.Error("func(LoginUser): No such user exists: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No such user exist",
		})
		return
	}

	if user.Password != req.Password {
		logger.Log.Error("func(LoginUser): Credentials do not match: ", err)
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Unauthorized Access",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.Email)
	if err != nil {
		logger.Log.Error("func(LoginUser): Unable to generate a new token: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to generate a token",
		})
		return
	}

	logger.Log.Infof("func(LoginUser): A new user logged in with id: %v ", user.ID)
	ctx.JSON(http.StatusAccepted, gin.H{
		"user-id":    user.ID,
		"user-name":  user.Name,
		"user-email": user.Email,
		"token":      token,
	})
}

func (u *UserController) UserActionOnStock(ctx *gin.Context) {
	var req models.RequestBodyTransaction

	userId, err := utils.ParseUserId(ctx)
	if err != nil {
		logger.Log.Error("func(UserActionOnStock): Cannot find user id from the token ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to fetch the user id from token",
		})
		return
	}

	stockIdParam := ctx.Param("stock_id")
	stockId, err := uuid.Parse(stockIdParam)
	if err != nil {
		logger.Log.Error("func(UserActionOnStock): Unable to parse the stock id: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to parse the user id",
		})
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Error("func(UserActionOnStock): Unable to parse the request body: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please validate the request body",
		})
		return
	}

	// if user sells the stock, logic to see if user has enough stocks to sell
	if req.Type == "sell" {
		logger.Log.Info(userId)
		logger.Log.Info(stockId)

		portfolio, err := u.PortfolioService.FetchUserPortfolioByStockId(ctx, userId, stockId)
		if err != nil {
			logger.Log.Error("func(UserActionOnStock): User do not have holding of this stock: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Use do not have holdings of this stock",
			})
			return
		}

		quantity, err := strconv.ParseFloat(portfolio.TotalQuantity, 64)
		if err != nil {
			logger.Log.Error("func(UserActionOnStock): Cannot parse the float: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot parse the float value for quantity",
			})
			return
		}

		if req.Quantity > quantity {
			logger.Log.Error("func(UserActionOnStock): User do not have enough stock to sell: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "You do not have enough stock to sell",
			})
			return
		}
	}

	stock, err := u.StockService.GetStockById(ctx, stockId)
	if err != nil {
		logger.Log.Error("func(UserActionOnStock): Unable to sretrieve the stock details: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("Cannot retrive details about stock : %v", stockId),
		})
		return
	}

	stockValuation, err := strconv.ParseFloat(stock.Valuation, 64)
	if err != nil {
		logger.Log.Error("func(UserActionOnStock): Unable to parse the string value to float: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprint("Unable to parse the string value to float: ", err),
		})
		return
	}
	price := stockValuation * req.Quantity

	transactionObject := models.TransactionRequestObject{
		UserID:          userId,
		StockId:         stockId,
		Quantity:        fmt.Sprintf("%.6f", req.Quantity),
		Type:            req.Type,
		Price:           fmt.Sprintf("%.4f", price),
		TransactionType: "transact",
	}

	transaction, err := u.TransactionService.RegisterTransaction(ctx, &transactionObject)
	if err != nil {
		logger.Log.Error("func(UserActionOnStock): Cannot complet the transaction: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to register the new transaction",
		})
		return
	}

	portfolioRequest := models.RecordPortfolioRequest{
		UserId:   userId,
		StockId:  stockId,
		Quantity: transaction.Quantity,
		Type:     transaction.Type,
	}

	_, err = u.PortfolioService.UpdateUserPortfolio(ctx, &portfolioRequest)
	if err != nil {
		logger.Log.Error("func(UserActionOnStock): Cannot record the portfolio: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to update the users portfolio",
		})
		return
	}

	logger.Log.Info("func(UserActionOnStock): Allocation of the stock successful")
	ctx.JSON(http.StatusAccepted, gin.H{
		"data": transaction,
	})
}

func (u *UserController) FetchUsersTransactionForToday(ctx *gin.Context) {
	userIdParam := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdParam)
	if err != nil {
		logger.Log.Error("func(FetchUsersTransactionForToday): Unable to parse the user id: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please validate the api url",
		})
		return
	}

	todaysDate := time.Now()
	logger.Log.Info(todaysDate)

	transactions, err := u.TransactionService.GetUserTransactionForToday(ctx, userId, todaysDate)
	if err != nil {
		logger.Log.Errorf("func(FetchUsersTransactionForToday): Cannot retrieve today's transaction for user with id: %v: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot retrive users transaction for today",
		})
		return
	}

	transactionResponse, err := utils.MapTransactions(ctx, "transact", u.StockService, transactions)
	if err != nil {
		logger.Log.Error("func(FetchUsersTransactionForToday): Unable to map the users transaction")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to map the users transaction",
		})
		return
	}

	logger.Log.Infof("func(FetchUsersTransactionForToday): Records for userid: %v fetched for today", userId)
	ctx.JSON(http.StatusAccepted, gin.H{
		"data": transactionResponse,
	})
}

func (u *UserController) FetchUsersRewardForToday(ctx *gin.Context) {
	userIdParam := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdParam)
	if err != nil {
		logger.Log.Error("func(FetchUsersRewardForToday): Unable to parse the user id: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please validate the api url",
		})
		return
	}

	todaysDate := time.Now()

	transactions, err := u.TransactionService.GetUserTransactionForToday(ctx, userId, todaysDate)
	if err != nil {
		logger.Log.Errorf("func(FetchUsersRewardForToday): Cannot retrieve today's transaction for user-id: %v: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot retrive users transaction for today",
		})
		return
	}

	transactionResponse, err := utils.MapTransactions(ctx, "transact", u.StockService, transactions)
	if err != nil {
		logger.Log.Error("func(FetchUsersRewardForToday): Unable to map the users transaction")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to map the users transaction",
		})
		return
	}

	logger.Log.Infof("func(FetchUsersRewardForToday): Records for userid: %v fetched for today", userId)
	ctx.JSON(http.StatusAccepted, gin.H{
		"data": transactionResponse,
	})
}

func (u *UserController) FetchUserPortfolio(ctx *gin.Context) {
	userId, err := utils.ParseUserId(ctx)
	if err != nil {
		logger.Log.Error("func(UserPortfolio): Cannot fetch user Id from the token: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to fetch the user id from token",
		})
		return
	}

	userPortfolio, err := u.PortfolioService.GetPortfolioByUserId(ctx, userId)
	if err != nil {
		logger.Log.Error("func(UserPortfolio) Cannot retrieve user's portfolio: ", err)
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "Unable to fetch the users portfolio",
		})
		return
	}

	logger.Log.Info("Successfully retireved users portfolio")
	ctx.JSON(http.StatusOK, gin.H{
		"data": utils.MapPortfolio(ctx, u.StockService, userPortfolio),
	})
}
