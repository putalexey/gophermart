package handlers

import (
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"go.uber.org/zap"
)

type Handlers struct {
	Logger *zap.SugaredLogger
	JWT    *ginjwt.GinJWTMiddleware
}

func New(logger *zap.SugaredLogger, jwtMiddleware *ginjwt.GinJWTMiddleware) *Handlers {
	return &Handlers{
		Logger: logger,
		JWT:    jwtMiddleware,
	}
}
