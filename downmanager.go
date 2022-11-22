package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

func getLink(c *gin.Context) {
	resp, err := http.Get("https://1fichier.com/?ecrarnm5mdj3ig863rvn")
	if err != nil {
		fmt.Errorf("Error while getting the link", err)
	}
	defer resp.Body.Close()
	body, err := httputil.DumpResponse(resp, true)
	fmt.Println(string(body))
	c.String(200, "Hello Man!")
}

func main() {
	router := gin.Default()
	router.GET("/link", getLink)

	router.Run("localhost:8080")
}
