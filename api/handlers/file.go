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

	adbkit.New("127.0.0.1", 5037).Select("emulator-5554").Push(file, "/sdcard/a.png")
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
