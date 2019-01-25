package router

import (
	"adbs/api/handlers"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	gin.SetMode("debug")

	api := r.Group("/api")
	{
		// 设备列表管理
		devices := api.Group("/devices")
		{
			// 获取设备列表
			devices.GET("/", handlers.GetDevices)
			// 连接设备
			devices.POST("/connect", handlers.ConnectDevice)
			// 断开设备
			devices.POST("/disconnect", handlers.DisconnectDevice)
		}

		// 单台设备管理
		device := api.Group("/device")
		{
			// 获取包列表
			device.GET("/packages", handlers.GetPackages)
			// 获取截屏
			device.GET("/screencap", handlers.ScreenCap)
			// 上传文件
			device.POST("/push", handlers.Push)
			// 拉取文件
			device.GET("/pull", handlers.Pull)
			// 获取目录
			device.GET("/dir", handlers.Dir)
			device.GET("/stat", handlers.Stat)
			// 模拟输入
			device.POST("/input", handlers.Input)

			// 处理websocket
			device.GET("/shell/ws", func(c *gin.Context) {
				handlers.WsHandler(c.Writer, c.Request)
			})
		}
	}

	//r.LoadHTMLGlob("templates/*")
	r.LoadHTMLFiles("templates/shell.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "shell.html", gin.H{
			"title": "Main website",
		})
	})

	r.Static("/static", "static")

	return r
}
