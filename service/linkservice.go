package service

import (
	"log"
	"net/http"
	"time"
	"zemoa/downmanager/database/config"
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
	InError   bool
	ErrorMsg  *string
	Length    uint32
	Filename  string
}

type LinkService struct {
	Db              *gorm.DB
	DownloadService *DownloadService
}

func (ls *LinkService) CreateLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLink := c.Query("link")
		log.Printf("Will create link with %s", paramLink)
		linkEntity := link.Create(paramLink, ls.Db)
		filename, rangeSupported, length, error := ls.DownloadService.GetLinkDetails(linkEntity.Link)
		log.Printf("filename <%s>, range supported <%t>, length <%d>", filename, rangeSupported, length)
		if error != nil {
			linkEntity.InError = true
			errMsg := "The file is unavailable"
			linkEntity.ErrorMsg = &errMsg
		} else {
			linkEntity.Filename = filename
			linkEntity.Rangesupported = rangeSupported
			linkEntity.Length = uint32(length)
		}
		link.Update(linkEntity, ls.Db)
		linkDto := convertLinkToDto(linkEntity)
		c.IndentedJSON(http.StatusCreated, linkDto)
	}

}

func (ls *LinkService) GetAllLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		allLinks := link.GetAll(ls.Db)
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

func (ls *LinkService) DeleteLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLinkRef := c.Param("linkref")
		log.Printf("Will delete link %s", paramLinkRef)
		link.DeleteByRef(uuid.MustParse(paramLinkRef), ls.Db)
		c.Writer.WriteHeader(http.StatusNoContent)
	}
}

func (ls *LinkService) StartDownloadLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLinkRef := c.Param("linkref")
		log.Printf("Start downloading %s", paramLinkRef)
		linkObj := link.GetByRef(uuid.MustParse(paramLinkRef), ls.Db)
		config := config.Get(ls.Db)
		ls.DownloadService.DownloadFile(linkObj.Link, linkObj.Filename, config.DownloadDir)
		startOrStopDownload(linkObj, true, ls.Db, c)
	}
}

func (ls *LinkService) StopDownloadLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLinkRef := c.Param("linkref")
		log.Printf("Stop downloading %s", paramLinkRef)
		linkObj := link.GetByRef(uuid.MustParse(paramLinkRef), ls.Db)
		startOrStopDownload(linkObj, false, ls.Db, c)
	}
}

func startOrStopDownload(linkObj *link.Link, start bool, db *gorm.DB, c *gin.Context) {
	if linkObj == nil {
		c.Writer.WriteHeader(http.StatusNotFound)
	} else {
		linkObj.Running = start
		link.Update(linkObj, db)
		c.IndentedJSON(http.StatusOK, convertLinkToDto(linkObj))
	}
}
func convertLinkToDto(link *link.Link) LinkDto {
	return LinkDto{
		Ref:       link.Ref,
		Link:      link.Link,
		Running:   link.Running,
		InError:   link.InError,
		CreatedAt: link.CreatedAt,
		ErrorMsg:  link.ErrorMsg,
		Length:    link.Length,
		Filename:  link.Filename,
	}
}
