package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"zemoa/downmanager/database/config"
	"zemoa/downmanager/database/link"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LinkDto struct {
	Ref        uuid.UUID
	CreatedAt  time.Time
	Link       string
	Running    bool
	InError    bool
	ErrorMsg   *string
	Size       uint64
	Downloaded uint64
	Filename   string
}

type LinkService struct {
	LinkRepo         *link.LinkRepo
	ConfigRepo       *config.ConfigRepo
	DownloadService  *DownloadService
	WebSocketService *WebSocketService
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
			linkEntity.Size = uint64(length)
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
		var listener = downloadListenerImpl{link: linkObj, linkRepo: ls.LinkRepo, websocketService: ls.WebSocketService}
		go ls.DownloadService.DownloadFile(linkObj.Link, linkObj.Filename, config.DownloadDir, &listener)
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
		Ref:        link.Ref,
		Link:       link.Link,
		Running:    link.Running,
		InError:    link.InError,
		CreatedAt:  link.CreatedAt,
		ErrorMsg:   link.ErrorMsg,
		Size:       link.Size,
		Downloaded: link.Downloaded,
		Filename:   link.Filename,
	}
}

type downloadListenerImpl struct {
	link              *link.Link
	linkRepo          *link.LinkRepo
	websocketService  *WebSocketService
	updateTriggerTime *time.Time
}

type downloadMessage struct {
	Linkref    uuid.UUID
	Finished   bool
	InError    bool
	ErrorMsg   string
	Total      uint64
	Downloaded uint64
}

func (dli *downloadListenerImpl) progress(passedByte uint64) {
	dli.linkRepo.UpdateDownloaded(dli.link.Ref, passedByte)
	now := time.Now()
	if dli.updateTriggerTime == nil || dli.updateTriggerTime.Before(now) {
		tmpTrigger := now.Add(1 * time.Second)
		dli.updateTriggerTime = &tmpTrigger
		progressMessage := downloadMessage{
			Linkref:    dli.link.Ref,
			Finished:   false,
			InError:    false,
			Total:      dli.link.Size,
			ErrorMsg:   "",
			Downloaded: passedByte,
		}
		dli.websocketService.BroadcastMessage(marshalDownloadMessage(&progressMessage))
	}

	if dli.link.Size == passedByte {
		link := dli.linkRepo.GetByRef(dli.link.Ref)
		link.Running = false
		dli.linkRepo.Update(link)
		endMessage := downloadMessage{
			Linkref:    dli.link.Ref,
			Finished:   true,
			InError:    false,
			Total:      dli.link.Size,
			ErrorMsg:   "",
			Downloaded: passedByte,
		}
		dli.websocketService.BroadcastMessage(marshalDownloadMessage(&endMessage))
		log.Printf("<%s> Download finish", dli.link.Ref)
	}
}

func marshalDownloadMessage(downloadMessage *downloadMessage) []byte {
	json, _ := json.Marshal(downloadMessage)
	return json
}
