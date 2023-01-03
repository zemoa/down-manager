package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"zemoa/downmanager/database"
	"zemoa/downmanager/database/config"
	"zemoa/downmanager/database/link"
	"zemoa/downmanager/service"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getLink(c *gin.Context) {
	resp, err := http.Get("https://1fichier.com/?ecrarnm5mdj3ig863rvn")
	if err != nil {
		log.Fatalf("Error while getting the link %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := httputil.DumpResponse(resp, true)
	log.Println(string(body))
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
	linkRepo := &link.LinkRepo{Db: db}
	configRepo := &config.ConfigRepo{Db: db}
	websocketService := service.NewWebSocket()
	downloadService := new(service.DownloadService)
	linkService := &service.LinkService{
		LinkRepo:         linkRepo,
		ConfigRepo:       configRepo,
		WebSocketService: websocketService,
		DownloadService:  downloadService,
	}
	configService := &service.ConfigService{
		ConfigRepo: configRepo,
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/link1", getLink)

	router.GET("/ws", websocketService.WebSocket())
	linkRoutes := router.Group("/links")
	linkRoutes.POST("", linkService.CreateLink())
	linkRoutes.GET("", linkService.GetAllLink())
	linkRoutes.DELETE(":linkref", linkService.DeleteLink())
	linkRoutes.PUT(":linkref/start", linkService.StartDownloadLink())
	linkRoutes.PUT(":linkref/stop", linkService.StopDownloadLink())

	configRoutes := router.Group("/config")
	configRoutes.GET("", configService.GetConfig())
	configRoutes.PATCH("", configService.UpdateConfig())

	router.Run("localhost:8080")
}
