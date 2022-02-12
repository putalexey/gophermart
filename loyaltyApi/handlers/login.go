package handlers

import (
	"errors"
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/putalexey/gophermart/loyaltyApi/repository"
	"github.com/putalexey/gophermart/loyaltyApi/requests"
	"github.com/putalexey/gophermart/loyaltyApi/utils"
	"go.uber.org/zap"
)

func Login(logger *zap.SugaredLogger, jwtMiddleware *ginjwt.GinJWTMiddleware, repo *repository.Repo) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		jwtMiddleware.LoginHandler(c)
	}
}
func Authenticator(repo repository.UserRepository) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		loginRequest := requests.LoginRequest{}
		err := c.ShouldBindJSON(&loginRequest)
		if err != nil {
			return nil, errors.New("missing Login or Password")
		}

		user, err := repo.FindUserByLogin(c, loginRequest.Login)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return nil, ginjwt.ErrFailedAuthentication
			} else {
				return nil, err
			}
		}

		if !utils.PasswordCheck(loginRequest.Password, user.Password) {
			return nil, ginjwt.ErrFailedAuthentication
		}
		return user, nil
	}
}
