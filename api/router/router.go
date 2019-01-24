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
		// 获取设备列表
		api.GET("/devices", handlers.GetDevices)
		// 连接设备
		api.POST("/connect", handlers.ConnectDevice)
		// 断开设备
		api.POST("/disconnect", handlers.DisconnectDevice)

		// 获取包列表
		api.GET("/device/packages", handlers.GetPackages)
		// 获取截屏
		api.GET("/device/screencap", handlers.ScreenCap)
		// 上传文件
		api.POST("/device/push", handlers.Push)
		// 拉取文件
		api.GET("/device/pull", handlers.Pull)
	}

	// 处理websocket
	r.GET("/ws/shell", func(c *gin.Context) {
		handlers.WsHandler(c.Writer, c.Request)
	})

	r.POST("/upload", handlers.Upload)

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
