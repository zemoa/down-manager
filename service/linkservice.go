package service

import (
	"net/http"
	"zemoa/downmanager/database/link"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		alink := c.Param("link")
		link := link.Create(*&alink, db)
		c.IndentedJSON(http.StatusCreated, link)
	}

}

func GetAllLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		allLinks := link.GetAll(db)
		c.IndentedJSON(http.StatusOK, allLinks)
	}
}
