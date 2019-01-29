package handlers

import (
	"adbs/adbkit"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"path"
)

const TEMP_PATH = "/data/local/tmp"

// Push 上传文件到 设备
func Push(c *gin.Context) {
	file, _ := c.FormFile("file")
	serial := c.Query("serial")
	p := c.Query("path")

	adbkit.New(CLENT_IP, CLENT_PORT).Select(serial).Push(file, p)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// Pull 从设备路径下载文件
func Pull(c *gin.Context) {
	serial := c.Query("serial")
	p := c.Query("path")

	content, err := adbkit.New(CLENT_IP, CLENT_PORT).Select(serial).Pull(p)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"message": err.Error()})
		return
	}
	c.Header("content-disposition", `attachment; filename=`+path.Base(p))
	c.Data(200, mime.TypeByExtension(path.Ext(p)), content)
}

func Install(c *gin.Context) {
	file, _ := c.FormFile("file")
	serial := c.Query("serial")
	p := TEMP_PATH + "/_install.apk"
	err := adbkit.New(CLENT_IP, CLENT_PORT).Select(serial).Push(file, p)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "upload file error"})
		return
	}
	_, err = adbkit.New(CLENT_IP, CLENT_PORT).Install(serial, p)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "apk install error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// Dir 获取某设备路径下得文件（夹）列表
func Dir(c *gin.Context) {
	serial := c.Query("serial")
	p := c.Query("path")

	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusGatewayTimeout, gin.H{"message": err})
			return
		}
	}()

	stats, err := adbkit.New(CLENT_IP, CLENT_PORT).Select(serial).Dir(p)
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

// Stat 获取某设备某路径文件详情
func Stat(c *gin.Context) {
	serial := c.Query("serial")
	p := c.Query("path")

	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusGatewayTimeout, gin.H{"message": err})
			return
		}
	}()

	stat, err := adbkit.New(CLENT_IP, CLENT_PORT).Select(serial).Stat(p)
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
