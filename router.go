package main

import (
	"github.com/Creedowl/NiuwaBI/handlers"
	"github.com/Creedowl/NiuwaBI/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func InitApp() *gin.Engine {
	logrus.Infoln("init app")
	authMiddleware, err := InitAuth()
	if err != nil {
		logrus.Fatalf("failed to init auth: %v", err)
	}
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	r.Use(utils.CustomLogger())
	gin.ForceConsoleColor()
	api := r.Group("/api")

	api.GET("/ping", utils.AutoWrap(handlers.Ping))
	api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/register", utils.AutoWrap(handlers.Register))
	api.GET("/refresh_token", authMiddleware.RefreshHandler)
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.POST("/test", utils.AutoWrap(handlers.Test))
		api.POST("/get_user_info", utils.AutoWrap(handlers.GetUserStatisticsInfo))
		workspace := api.Group("/workspace")
		{
			workspace.POST("/get", utils.AutoWrap(handlers.GetWorkspace))
			workspace.POST("/get_all", utils.AutoWrap(handlers.GetAllWorkspaces))
			workspace.POST("/create", utils.AutoWrap(handlers.CreateWorkspace))
			workspace.POST("/update", utils.AutoWrap(handlers.UpdateWorkspace))
			workspace.POST("/remove", utils.AutoWrap(handlers.RemoveWorkspace))
			workspace.POST("/test_conn", utils.AutoWrap(handlers.TestConn))
		}

		report := api.Group("/report")
		{
			report.POST("/get", utils.AutoWrap(handlers.GetReport))
			report.POST("/get_all", utils.AutoWrap(handlers.GetAllReports))
			report.POST("/create", utils.AutoWrap(handlers.CreateReport))
			report.POST("/update", utils.AutoWrap(handlers.UpdateReport))
			report.POST("/remove", utils.AutoWrap(handlers.RemoveReport))
			report.POST("/execute", utils.AutoWrap(handlers.ExecuteReport))
		}
	}

	return r
}
