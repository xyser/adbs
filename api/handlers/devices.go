package handlers

import (
	"adbs/shell"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// 获取多个项目
func GetDevices(c *gin.Context) {
	// 设备列表
	devices, err := shell.Lists()
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

func ConnectDevice(c *gin.Context) {
	var message = "success"

	var post struct {
		Ip string `json:"ip"`
	}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	bo, err := shell.Connect(post.Ip)
	if err != nil || !bo {
		message = fmt.Sprintf("devices connect: %s", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func DisconnectDevice(c *gin.Context) {
	var message = "success"

	var post struct {
		Ip string `json:"ip"`
	}
	if err := c.ShouldBindJSON(&post); err != nil {
		post.Ip = "all"
	}
	bo, err := shell.Disconnect(post.Ip)
	if err != nil || !bo {
		message = fmt.Sprintf("devices connect: %s", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func ScreenCap(c *gin.Context) {
	buffer, err := shell.Screencap()
	if err == nil {
		c.Writer.Header().Set("Content-Type", "image/png")
		c.Writer.Header().Set("Content-Length", strconv.Itoa(len(buffer)))
		if _, err := c.Writer.Write(buffer); err != nil {
			log.Println("unable to write image.")
		}
	} else {
		c.String(http.StatusOK, err.Error())
	}

}
