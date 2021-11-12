package main

import (
	"github.com/Creedowl/NiuwaBI/handlers"
	"github.com/Creedowl/NiuwaBI/utils"
	"github.com/gin-gonic/gin"
)

func InitApp() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/ping", utils.AutoWrap(handlers.Ping))
		api.POST("/test", utils.AutoWrap(handlers.Test))
	}

	return r
}
