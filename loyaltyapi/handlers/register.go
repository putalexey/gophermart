package handlers

import (
	"errors"
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/putalexey/gophermart/loyaltyapi/models"
	"github.com/putalexey/gophermart/loyaltyapi/repository"
	"github.com/putalexey/gophermart/loyaltyapi/requests"
	"github.com/putalexey/gophermart/loyaltyapi/responses"
	"github.com/putalexey/gophermart/loyaltyapi/utils"
	"net/http"
)

var ErrLoginUsed = errors.New("login used")

func Register(mw *ginjwt.GinJWTMiddleware, repo repository.UserRepository) func(c *gin.Context) {
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

		user := &models.User{
			UUID:     uuid.NewString(),
			Login:    registerRequest.Login,
			Password: utils.PasswordHash(registerRequest.Password),
		}

		_, err = repo.CreateUser(c, user)
		if err != nil {
			responses.JSONError(c, http.StatusBadRequest, err)
			return
		}

		tokenString, expire, err := mw.TokenGenerator(user)
		if err != nil {
			responses.JSONError(c, http.StatusInternalServerError, err)
			return
		}

		// set cookie
		if mw.SendCookie {
			expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
			maxage := int(expireCookie.Unix() - mw.TimeFunc().Unix())

			if mw.CookieSameSite != 0 {
				c.SetSameSite(mw.CookieSameSite)
			}

			c.SetCookie(
				mw.CookieName,
				tokenString,
				maxage,
				"/",
				mw.CookieDomain,
				mw.SecureCookie,
				mw.CookieHTTPOnly,
			)
		}

		mw.LoginResponse(c, http.StatusOK, tokenString, expire)
	}
}
