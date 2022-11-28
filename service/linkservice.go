package service

import (
	"net/http"
	"time"
	"zemoa/downmanager/database/link"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LinkDto struct {
	ID        uint
	CreatedAt time.Time
	link      string
	running   bool
	inerror   bool
}

func CreateLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLink := c.Param("link")
		linkEntity := link.Create(*&paramLink, db)
		linkDto := LinkDto{
			ID:        linkEntity.ID,
			link:      linkEntity.Link,
			running:   linkEntity.Running,
			inerror:   linkEntity.InError,
			CreatedAt: linkEntity.CreatedAt,
		}
		c.IndentedJSON(http.StatusCreated, linkDto)
	}

}

func GetAllLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		allLinks := link.GetAll(db)
		c.IndentedJSON(http.StatusOK, allLinks)
	}
}
