package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/putalexey/gophermart/loyaltyapi/responses"
	"net/http"
)

func (h *Handlers) NotImplemented(c *gin.Context) {
	responses.JSONError(c, http.StatusNotImplemented, "not implemented")
}
