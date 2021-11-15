package main

import (
	"github.com/Creedowl/NiuwaBI/handlers"
	"github.com/Creedowl/NiuwaBI/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitApp() *gin.Engine {
	logrus.Infoln("init app")
	authMiddleware, err := InitAuth()
	if err != nil {
		logrus.Fatalf("failed to init auth: %v", err)
	}
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/ping", utils.AutoWrap(handlers.Ping))
	api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/register", utils.AutoWrap(handlers.Register))
	api.GET("/refresh_token", authMiddleware.RefreshHandler)
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.POST("/test", utils.AutoWrap(handlers.Test))
	}

	return r
}
