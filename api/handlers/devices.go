package handlers

import (
	"adbs/adbkit"
	"adbs/shell"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"strconv"
)

// 获取多个项目
func GetDevices(c *gin.Context) {
	// 设备列表
	devices, err := adbkit.New("127.0.0.1", 5037).Lists()
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

	ip := c.PostForm("ip")

	if ip == "" || net.ParseIP(ip) == nil {
		c.JSON(http.StatusOK, gin.H{"message": "IP Error"})
		return
	}

	bo, err := shell.Connect(ip)
	if err != nil || !bo {
		message = fmt.Sprintf("devices connect: %s", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func DisconnectDevice(c *gin.Context) {
	var message = "success"

	ip := c.PostForm("ip")
	if ip == "" || net.ParseIP(ip) == nil {
		c.JSON(http.StatusOK, gin.H{"message": "IP Error"})
		return
	}

	bo, err := shell.Disconnect(ip)
	if err != nil || !bo {
		message = fmt.Sprintf("devices connect: %s", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func Screencap(c *gin.Context) {
	serial := c.Query("serial")
	buffer, err := adbkit.New("127.0.0.1", 5037).Screencap(serial)
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

func WindowSize(c *gin.Context) {
	serial := c.Query("serial")
	w, h, err := adbkit.New("127.0.0.1", 5037).ScreenSize(serial)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"width": w, "height": h})
	} else {
		c.String(http.StatusOK, err.Error())
	}

}
