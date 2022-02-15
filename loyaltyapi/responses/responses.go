package responses

import (
	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, code int, err string) {
	c.JSON(code, ErrorResponse{
		Error: err,
	})
}

func JSON(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
