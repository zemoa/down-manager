package service

import (
	"log"
	"net/http"
	"time"
	"zemoa/downmanager/database/config"
	"zemoa/downmanager/database/link"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	LinkRepo        *link.LinkRepo
	ConfigRepo      *config.ConfigRepo
	DownloadService *DownloadService
}

func (ls *LinkService) CreateLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLink := c.Query("link")
		log.Printf("Will create link with %s", paramLink)
		linkEntity := ls.LinkRepo.Create(paramLink)
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
		ls.LinkRepo.Update(linkEntity)
		linkDto := convertLinkToDto(linkEntity)
		c.IndentedJSON(http.StatusCreated, linkDto)
	}

}

func (ls *LinkService) GetAllLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		allLinks := ls.LinkRepo.GetAll()
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
		ls.LinkRepo.DeleteByRef(uuid.MustParse(paramLinkRef))
		c.Writer.WriteHeader(http.StatusNoContent)
	}
}

func (ls *LinkService) StartDownloadLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLinkRef := c.Param("linkref")
		log.Printf("Start downloading %s", paramLinkRef)
		linkObj := ls.LinkRepo.GetByRef(uuid.MustParse(paramLinkRef))
		config := ls.ConfigRepo.Get()
		ls.DownloadService.DownloadFile(linkObj.Link, linkObj.Filename, config.DownloadDir)
		ls.startOrStopDownload(linkObj, true, c)
	}
}

func (ls *LinkService) StopDownloadLink() func(c *gin.Context) {
	return func(c *gin.Context) {
		paramLinkRef := c.Param("linkref")
		log.Printf("Stop downloading %s", paramLinkRef)
		linkObj := ls.LinkRepo.GetByRef(uuid.MustParse(paramLinkRef))
		ls.startOrStopDownload(linkObj, false, c)
	}
}

func (ls *LinkService) startOrStopDownload(linkObj *link.Link, start bool, c *gin.Context) {
	if linkObj == nil {
		c.Writer.WriteHeader(http.StatusNotFound)
	} else {
		linkObj.Running = start
		ls.LinkRepo.Update(linkObj)
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
