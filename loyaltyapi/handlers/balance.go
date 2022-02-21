package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/putalexey/gophermart/loyaltyapi/models"
	"github.com/putalexey/gophermart/loyaltyapi/repository"
	"github.com/putalexey/gophermart/loyaltyapi/requests"
	"github.com/putalexey/gophermart/loyaltyapi/responses"
	"github.com/putalexey/gophermart/loyaltyapi/utils"
	"net/http"
	"time"
)

func (h *Handlers) GetUserBalance(repo repository.BalanceRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		tmpUser, exists := c.Get(models.UserIdentityKey)
		if !exists {
			h.Logger.Error("user identity is not in context. Forget to add jwtMiddleware.MiddlewareFunc() to the router group?")
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}
		user, ok := tmpUser.(*models.User)
		if !ok {
			h.Logger.Error("loaded identity is not models.User", tmpUser)
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}

		balance, err := repo.GetUserBalance(c, user.UUID)
		if err != nil {
			h.Logger.Error(err)
			responses.JSONError(c, http.StatusInternalServerError, err.Error())
			return
		}

		responses.JSON(c, http.StatusOK, balance)
	}
}

func (h *Handlers) GetBalanceWithdrawals(repo repository.BalanceRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		tmpUser, exists := c.Get(models.UserIdentityKey)
		if !exists {
			h.Logger.Error("user identity is not in context. Forget to add jwtMiddleware.MiddlewareFunc() to the router group?")
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}
		user, ok := tmpUser.(*models.User)
		if !ok {
			h.Logger.Error("loaded identity is not models.User", tmpUser)
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}

		withdrawals, err := repo.GetBalanceWithdrawals(c, user)
		if err != nil {
			h.Logger.Error(err)
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}

		responses.JSON(c, http.StatusOK, withdrawals)
	}
}

func (h *Handlers) BalanceWithdraw(repo repository.BalanceRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		tmpUser, exists := c.Get(models.UserIdentityKey)
		if !exists {
			h.Logger.Error("user identity is not in context. Forget to add jwtMiddleware.MiddlewareFunc() to the router group?")
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}
		user, ok := tmpUser.(*models.User)
		if !ok {
			h.Logger.Error("loaded identity is not models.User", tmpUser)
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}

		var request requests.BalanceWithdraw
		err := c.ShouldBindJSON(&request)
		if err != nil {
			responses.JSONError(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		if !utils.CheckOrderNumber(request.Order) {
			responses.JSONError(c, http.StatusUnprocessableEntity, "wrong order number")
			return
		}

		withdrawal := &models.Withdrawal{
			UUID:        uuid.NewString(),
			UserUUID:    user.UUID,
			Order:       request.Order,
			Sum:         request.Sum,
			ProcessedAt: time.Now(),
		}
		_, err = repo.BalanceWithdraw(c, withdrawal)
		if err != nil {
			if errors.Is(err, repository.ErrNotEnoughBalance) {
				responses.JSONError(c, http.StatusPaymentRequired, err.Error())
			} else {
				h.Logger.Error(err)
				responses.JSONError(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		responses.JSON(c, http.StatusOK, nil)
	}
}
