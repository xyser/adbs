package handlers

import (
	"adbs/adbkit"
	"adbs/shell/packages"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取多个项目
// TODO:: 尚未支持切换设备
func GetPackages(c *gin.Context) {
	// 设备列表
	devices, err := packages.List()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("devices error: %s", err.Error()),
		})
	} else {
		if len(devices) == 0 {
			c.JSON(http.StatusOK, make([]string, 0))
		} else {
			c.JSON(http.StatusOK, devices)
		}
	}
}

// ClearPackage 用于 清理包缓存
func ClearPackage(c *gin.Context) {
	pkg := c.PostForm("package")
	serial := c.Query("serial")

	if err := adbkit.New(CLENT_IP, CLENT_PORT).PackageClear(serial, pkg); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("error: %s", err.Error())})
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
