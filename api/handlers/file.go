package handlers

import (
	"adbs/adbkit"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"path"
)

// 上传文件到 设备
func Push(c *gin.Context) {
	file, _ := c.FormFile("file")
	serial := c.Query("serial")
	p := c.Query("path")

	adbkit.New("127.0.0.1", 5037).Select(serial).Push(file, p)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 从设备路径下载文件
func Pull(c *gin.Context) {
	serial := c.Query("serial")
	p := c.Query("path")

	content, err := adbkit.New("127.0.0.1", 5037).Select(serial).Pull(p)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"message": err.Error()})
		return
	}
	c.Header("content-disposition", `attachment; filename=`+path.Base(p))
	c.Data(200, mime.TypeByExtension(path.Ext(p)), content)
}

func Dir(c *gin.Context) {
	serial := c.Query("serial")
	p := c.Query("path")

	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusGatewayTimeout, gin.H{"message": err})
			return
		}
	}()

	stats, err := adbkit.New("127.0.0.1", 5037).Select(serial).Dir(p)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"message": err.Error()})
		return
	}
	var resp []gin.H
	for _, stat := range stats {
		var s = gin.H{
			"name":     stat.Name,
			"size":     stat.Size,
			"mode":     stat.Mode.String(),
			"mod_time": stat.ModTime,
		}
		resp = append(resp, s)
	}

	c.JSON(http.StatusOK, resp)
}

func Stat(c *gin.Context) {
	serial := c.Query("serial")
	p := c.Query("path")

	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusGatewayTimeout, gin.H{"message": err})
			return
		}
	}()

	stat, err := adbkit.New("127.0.0.1", 5037).Select(serial).Stat(p)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"message": err.Error()})
		return
	}
	var s = gin.H{
		"name":     stat.Name,
		"size":     stat.Size,
		"mode":     stat.Mode.String(),
		"mod_time": stat.ModTime,
	}

	c.JSON(http.StatusOK, s)
}
