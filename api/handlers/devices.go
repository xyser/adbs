package handlers

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"adbs/adbkit"
	"adbs/shell"

	"github.com/gin-gonic/gin"
)

const CLENT_IP = "127.0.0.1"
const CLENT_PORT = 5037

// 获取多个项目
func GetDevices(c *gin.Context) {
	// 设备列表
	devices, err := adbkit.New(CLENT_IP, CLENT_PORT).Lists()
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
	if ip == "" {
		c.JSON(http.StatusOK, gin.H{"message": "IP Empty"})
		return
	}
	ips := strings.Split(ip, ":")
	if net.ParseIP(ips[0]) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "IP Error"})
		return
	}

	var port int
	if len(ips) > 1 {
		port, _ = strconv.Atoi(ips[1])
		ip = ips[0]
	} else {
		port = 5555
	}
	bo, err := adbkit.New(CLENT_IP, CLENT_PORT).Connect(ip, port)
	if err != nil || !bo {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("devices connect: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func DisconnectDevice(c *gin.Context) {
	var message = "success"

	ip := c.PostForm("ip")
	if ip == "" {
		c.JSON(http.StatusOK, gin.H{"message": "IP Empty"})
		return
	}
	ips := strings.Split(ip, ":")
	if net.ParseIP(ips[0]) == nil {
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
	channel := c.Query("channel")

	var buffer []byte
	var err error
	if channel == "shell" {
		buffer, err = shell.Screencap(serial)
	} else {
		buffer, err = adbkit.New(CLENT_IP, CLENT_PORT).Screencap(serial)
	}

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
	w, h, err := adbkit.New(CLENT_IP, CLENT_PORT).ScreenSize(serial)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"width": w, "height": h})
	} else {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}

}
