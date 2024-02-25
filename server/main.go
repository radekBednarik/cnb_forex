package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello, world!"})
	})
	g.Run()
}
