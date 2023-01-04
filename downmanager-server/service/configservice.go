package service

import (
	"log"
	"net/http"
	"zemoa/downmanager/database/config"

	"github.com/gin-gonic/gin"
)

type ConfigDto struct {
	DownloadDir string
}

type ConfigService struct {
	ConfigRepo *config.ConfigRepo
}

func (cs *ConfigService) GetConfig() func(c *gin.Context) {
	return func(c *gin.Context) {
		configModle := cs.ConfigRepo.Get()
		c.IndentedJSON(http.StatusOK, ConfigDto{DownloadDir: configModle.DownloadDir})
	}
}

func (cs *ConfigService) UpdateConfig() func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Print("Will update config")
		configDto := new(ConfigDto)
		err := c.Bind(configDto)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			configModel := cs.ConfigRepo.Update(&config.ConfigModel{DownloadDir: &configDto.DownloadDir})
			c.IndentedJSON(http.StatusOK, ConfigDto{DownloadDir: *configModel.DownloadDir})
		}
	}
}
