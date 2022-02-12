package responses

import (
	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, code int, err error) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: err.Error(),
	})
}

func JSON(c *gin.Context, code int, data map[string]interface{}) {
	if _, ok := data["code"]; !ok {
		data["code"] = code
	}
	c.JSON(code, data)
}
