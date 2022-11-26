package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"zemoa/downmanager/database"
	"zemoa/downmanager/service"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func getLink(c *gin.Context) {
	resp, err := http.Get("https://1fichier.com/?ecrarnm5mdj3ig863rvn")
	if err != nil {
		log.Fatalf("Error while getting the link %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := httputil.DumpResponse(resp, true)
	fmt.Println(string(body))
	c.String(200, "Hello Man!")
}

type DLForm1Fichier struct {
	link string
	adz  string
}

func GiveFileLink1Fichier(content string) DLForm1Fichier {
	parsed, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err.Error())
	}
	form := parsed.Find("form").First()
	attrAction, _ := form.Attr("action")
	adz, _ := form.Find("input[name=adz]").First().Attr("value")

	return DLForm1Fichier{link: attrAction, adz: adz}
}

func findForm(content *goquery.Document) *goquery.Selection {
	return content.Find("form")
}

func main() {
	db := database.Init(".")
	router := gin.Default()
	router.GET("/link1", getLink)
	router.POST("/links", service.CreateLink(db))
	router.GET("/links", service.GetAllLink(db))

	router.Run("localhost:8080")
}
