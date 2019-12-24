package handlers

import (
	"net/http"

	"adbs/shell"

	"github.com/gin-gonic/gin"
)

func Input(c *gin.Context) {
	// serial := c.Query("serial")
	command := c.PostForm("command")
	arg := c.PostForm("arg")

	_, err := shell.Input(command, arg)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
