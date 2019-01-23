package handlers

import (
	"adbs/adbkit"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	adbkit.New("127.0.0.1", 5037).Select("351BBJPAUULH").Push(file, "/sdcard")
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
