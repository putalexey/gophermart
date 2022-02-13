package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/putalexey/gophermart/loyaltyApi/models"
	"github.com/putalexey/gophermart/loyaltyApi/responses"
	"net/http"
)

func CurrentUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		user, exists := c.Get(models.UserIdentityKey)

		responses.JSON(c, http.StatusOK,
			map[string]interface{}{
				"user":   user,
				"exists": exists,
			},
		)
	}
}
