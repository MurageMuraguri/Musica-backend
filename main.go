package main

import (
	"musica/model/migrate"

	"github.com/gin-gonic/gin"
)

func main() {
	migrate.Migrate()
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "kasee weh"})
	})

	r.Run(":9000")
}
