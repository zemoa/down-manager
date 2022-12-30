package service

import (
	"log"
	"net/http"
	"zemoa/downmanager/database/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigDto struct {
	DownloadDir string
}

type ConfigService struct {
	Db *gorm.DB
}

func (cs *ConfigService) GetConfig() func(c *gin.Context) {
	return func(c *gin.Context) {
		configModle := config.Get(cs.Db)
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
			configModel := config.Update(cs.Db, &config.ConfigModel{DownloadDir: &configDto.DownloadDir})
			c.IndentedJSON(http.StatusOK, ConfigDto{DownloadDir: *configModel.DownloadDir})
		}
	}
}
