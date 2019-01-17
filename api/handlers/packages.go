package handlers

import (
	"adbs/shell/packages"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取多个项目
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
