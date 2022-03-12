package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()


	log.Println(fingerprint)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "kasee weh"})
	})

	r.Run(":9000")
}
