package handlers

import (
	"adbs/adbkit"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"path"
)

func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")

	adbkit.New("127.0.0.1", 5037).Select("emulator-5554").Push(file, "/sdcard/a.png")
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

// 上传文件到 设备
func Push(c *gin.Context) {
	file, _ := c.FormFile("file")
	serial := c.Param("serial")
	p := c.Query("path")

	adbkit.New("127.0.0.1", 5037).Select(serial).Push(file, p)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 从设备路径下载文件
func Pull(c *gin.Context) {
	serial := c.Param("serial")
	p := c.Query("path")

	content, _ := adbkit.New("127.0.0.1", 5037).Select(serial).Pull(p)
	c.Header("content-disposition", `attachment; filename=`+path.Base(p))
	c.Data(200, mime.TypeByExtension(path.Ext(p)), content)
}
