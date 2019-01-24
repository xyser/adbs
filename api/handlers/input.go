package handlers

import (
	"adbs/shell"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Input(c *gin.Context) {
	// serial := c.Query("serial")
	command := c.Query("command")
	arg := c.Query("arg")

	_, err := shell.Input(command, arg)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
