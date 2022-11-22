package main

import "github.com/gin-gonic/gin"

func getLink(c *gin.Context) {
	c.String(200, "Hello Man!")
}

func main() {
	router := gin.Default()
	router.GET("/link", getLink)

	router.Run("localhost:8080")
}
