package router

import (
	"adbs/api/handlers"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode("debug")

	api := r.Group("/api")
	{
		//获取设备列表
		api.GET("/devices", handlers.GetDevices)
		// 连接设备
		api.POST("/devices/connect", handlers.ConnectDevice)
		// 断开设备
		api.POST("/devices/disconnect", handlers.DisconnectDevice)

		api.GET("/packages", handlers.GetPackages)
	}

	// 处理websocket
	r.GET("/ws", func(c *gin.Context) {
		handlers.WsHandler(c.Writer, c.Request)
	})

	return r
}
