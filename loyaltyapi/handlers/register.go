package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/putalexey/gophermart/loyaltyapi/models"
	"github.com/putalexey/gophermart/loyaltyapi/repository"
	"github.com/putalexey/gophermart/loyaltyapi/requests"
	"github.com/putalexey/gophermart/loyaltyapi/responses"
	"github.com/putalexey/gophermart/loyaltyapi/utils"
	"net/http"
)

func (h *Handlers) Register(repo repository.UserRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		registerRequest := requests.RegisterRequest{}
		err2 := c.ShouldBindJSON(&registerRequest)
		if err2 != nil {
			responses.JSONError(c, http.StatusBadRequest, err2.Error())
			return
		}

		if registerRequest.Login == "" || registerRequest.Password == "" {
			responses.JSONError(c, http.StatusBadRequest, "login or password is empty")
			return
		}

		_, err := repo.FindUserByLogin(c, registerRequest.Login)

		if !errors.Is(err, repository.ErrNotFound) {
			h.Logger.Error(err)
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}

		if err == nil {
			// user found
			responses.JSONError(c, http.StatusConflict, "login already in use")
			return
		}

		hash, err := utils.PasswordHash(registerRequest.Password)
		if err != nil {
			h.Logger.Error(err)
			responses.JSONError(c, http.StatusInternalServerError, "server error")
			return
		}

		user := &models.User{
			UUID:     uuid.NewString(),
			Login:    registerRequest.Login,
			Password: hash,
		}

		_, err = repo.CreateUser(c, user)
		if err != nil {
			responses.JSONError(c, http.StatusBadRequest, err.Error())
			return
		}

		tokenString, expire, err := h.JWT.TokenGenerator(user)
		if err != nil {
			responses.JSONError(c, http.StatusInternalServerError, err.Error())
			return
		}

		// set cookie
		if h.JWT.SendCookie {
			expireCookie := h.JWT.TimeFunc().Add(h.JWT.CookieMaxAge)
			maxage := int(expireCookie.Unix() - h.JWT.TimeFunc().Unix())

			if h.JWT.CookieSameSite != 0 {
				c.SetSameSite(h.JWT.CookieSameSite)
			}

			c.SetCookie(
				h.JWT.CookieName,
				tokenString,
				maxage,
				"/",
				h.JWT.CookieDomain,
				h.JWT.SecureCookie,
				h.JWT.CookieHTTPOnly,
			)
		}

		h.JWT.LoginResponse(c, http.StatusOK, tokenString, expire)
	}
}
