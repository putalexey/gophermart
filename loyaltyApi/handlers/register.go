package handlers

import (
	"errors"
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/putalexey/gophermart/loyaltyApi/repository"
	"github.com/putalexey/gophermart/loyaltyApi/requests"
	"github.com/putalexey/gophermart/loyaltyApi/responses"
	"go.uber.org/zap"
	"net/http"
)

var ErrLoginUsed = errors.New("login used")

func Register(logger *zap.SugaredLogger, middleware *ginjwt.GinJWTMiddleware, repo repository.UserRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		registerRequest := requests.RegisterRequest{}
		err2 := c.BindJSON(&registerRequest)
		if err2 != nil {
			responses.JSONError(c, http.StatusBadRequest, err2)
			return
		}

		_, err := repo.FindUserByLogin(c, registerRequest.Login)
		if err == nil {
			// user found
			responses.JSONError(c, http.StatusConflict, ErrLoginUsed)
			return
		}

		responses.JSONError(c, http.StatusInternalServerError, errors.New("not implemented"))

		//if user == nil || utils.PasswordCheck(loginRequest.Password, user.Password) {
		//	return nil, ginjwt.ErrFailedAuthentication
		//}
		//return user, nil
	}
}
