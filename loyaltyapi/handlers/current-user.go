package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/putalexey/gophermart/loyaltyapi/models"
	"github.com/putalexey/gophermart/loyaltyapi/responses"
	"net/http"
)

func (h *Handlers) CurrentUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		user, _ := c.Get(models.UserIdentityKey)
		responses.JSON(c, http.StatusOK, user)
	}
}
