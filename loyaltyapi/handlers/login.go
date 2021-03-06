package handlers

import (
	"errors"
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/putalexey/gophermart/loyaltyapi/repository"
	"github.com/putalexey/gophermart/loyaltyapi/requests"
	"github.com/putalexey/gophermart/loyaltyapi/utils"
)

func (h *Handlers) Authenticator(repo repository.UserRepository) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		loginRequest := requests.LoginRequest{}
		err := c.ShouldBindJSON(&loginRequest)
		if err != nil {
			return nil, errors.New("missing Login or Password")
		}

		if loginRequest.Login == "" || loginRequest.Password == "" {
			return nil, errors.New("login or password is empty")
		}

		user, err := repo.FindUserByLogin(c, loginRequest.Login)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
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
