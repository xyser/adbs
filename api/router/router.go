package router

import (
	"adbs/api/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 开启跨域
	r.Use(Cors())

	// 开启压缩
	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	gin.SetMode("debug")

	api := r.Group("/api")
	{
		// 设备列表管理
		devices := api.Group("/devices")
		{
			// 获取设备列表
			devices.GET("", handlers.GetDevices)
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
			device.GET("/screencap", handlers.Screencap)
			// 上传文件
			device.POST("/push", handlers.Push)
			// 拉取文件
			device.GET("/pull", handlers.Pull)
			// 获取目录
			device.GET("/dir", handlers.Dir)
			device.GET("/stat", handlers.Stat)
			// 模拟输入
			device.POST("/input", handlers.Input)
			// 获取屏幕尺寸
			device.GET("/window/size", handlers.WindowSize)

			// 处理websocket
			device.GET("/shell/ws", func(c *gin.Context) {
				handlers.WsHandler(c.Writer, c.Request)
			})
		}
	}

	//r.LoadHTMLGlob("templates/*")
	r.LoadHTMLFiles("templates/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	r.Static("/static", "static")

	return r
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//				允许跨域设置																										可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //	跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //	处理请求
	}
}
