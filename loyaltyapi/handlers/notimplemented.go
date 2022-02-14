package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/putalexey/gophermart/loyaltyapi/responses"
	"net/http"
)

func NotImplemented(c *gin.Context) {
	responses.JSONError(c, http.StatusNotImplemented, "not implemented")
}
