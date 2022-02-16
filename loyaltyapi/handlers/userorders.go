package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/putalexey/gophermart/loyaltyapi/models"
	"github.com/putalexey/gophermart/loyaltyapi/repository"
	"github.com/putalexey/gophermart/loyaltyapi/responses"
	"github.com/putalexey/gophermart/loyaltyapi/utils"
	"io"
	"net/http"
	"time"
)

func (h *Handlers) UserCreateOrder(orderRepo repository.OrderRepository, jobRepo repository.JobRepository) func(*gin.Context) {
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

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			responses.JSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		number := string(body)
		if valid := utils.CheckOrderNumber(number); !valid {
			responses.JSONError(c, http.StatusUnprocessableEntity, "wrong order number")
			return
		}

		existingOrder, err := orderRepo.GetOrderByNumber(c, number)
		if err == nil || !errors.Is(err, repository.ErrNotFound) {
			if err != nil {
				h.Logger.Error(err)
				responses.JSONError(c, http.StatusInternalServerError, "server error")
				return
			}
			if existingOrder.UserUUID != user.UUID {
				responses.JSONError(c, http.StatusConflict, "already loaded by another user")
				return
			}
			responses.JSON(c, http.StatusOK, nil)
			return
		}

		order := models.NewOrder()
		order.UserUUID = user.UUID
		order.Number = number

		_, err = orderRepo.CreateOrder(c, order)
		if err != nil {
			responses.JSONError(c, http.StatusInternalServerError, err.Error())
			return
		}

		job := &models.Job{
			UUID:      uuid.NewString(),
			OrderUUID: order.UUID,
			ProceedAt: time.Now(),
			Tries:     0,
		}
		err = jobRepo.CreateJob(c, job)
		if err != nil {
			h.Logger.Error("failed to create job", err)
		}

		responses.JSON(c, http.StatusAccepted, nil)
	}
}

func (h *Handlers) UserGetOrders(repo repository.OrderRepository) func(*gin.Context) {
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
		orders, err := repo.GetUserOrders(c, user)
		if err != nil {
			responses.JSONError(c, http.StatusInternalServerError, err.Error())
			return
		}

		responses.JSON(c, http.StatusOK, orders)
	}
}
