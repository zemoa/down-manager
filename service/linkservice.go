package service

import (
	"log"
	"net/http"
	"time"
	"zemoa/downmanager/database/link"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LinkDto struct {
	Ref       uuid.UUID
	CreatedAt time.Time
	Link      string
	Running   bool
	Inerror   bool
}

func CreateLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLink := c.Query("link")
		log.Printf("Will create link with %s", paramLink)
		linkEntity := link.Create(paramLink, db)
		linkDto := convertLinkToDto(linkEntity)
		c.IndentedJSON(http.StatusCreated, linkDto)
	}

}

func GetAllLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		allLinks := link.GetAll(db)
		var allLinksDto []LinkDto
		for _, link := range allLinks {
			allLinksDto = append(allLinksDto, convertLinkToDto(&link))
		}
		if allLinksDto == nil {
			allLinksDto = []LinkDto{}
		}
		c.IndentedJSON(http.StatusOK, allLinksDto)
	}
}

func DeleteLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLinkRef := c.Param("linkref")
		log.Printf("Will delete link %s", paramLinkRef)
		link.DeleteByRef(uuid.MustParse(paramLinkRef), db)
		c.Writer.WriteHeader(http.StatusNoContent)
	}
}

func StartDownloadLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLinkRef := c.Param("linkref")
		log.Printf("Start downloading %s", paramLinkRef)
		startOrStopDownload(paramLinkRef, true, db, c)
	}
}

func startOrStopDownload(paramLinkRef string, start bool, db *gorm.DB, c *gin.Context) {
	linkObj := link.GetByRef(uuid.MustParse(paramLinkRef), db)
	if linkObj == nil {
		c.Writer.WriteHeader(http.StatusNotFound)
	} else {
		linkObj.Running = start
		link.Update(linkObj, db)
		c.IndentedJSON(http.StatusOK, convertLinkToDto(linkObj))
	}
}

func StopDownloadLink(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLinkRef := c.Param("linkref")
		log.Printf("Stop downloading %s", paramLinkRef)
		startOrStopDownload(paramLinkRef, false, db, c)
	}
}

func convertLinkToDto(link *link.Link) LinkDto {
	return LinkDto{
		Ref:       link.Ref,
		Link:      link.Link,
		Running:   link.Running,
		Inerror:   link.InError,
		CreatedAt: link.CreatedAt,
	}
}
