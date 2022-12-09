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

func GetConfig(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		configModle := config.Get(db)
		c.IndentedJSON(http.StatusOK, ConfigDto{DownloadDir: configModle.DownloadDir})
	}
}

func UpdateConfig(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Print("Will update config")
		configDto := new(ConfigDto)
		err := c.Bind(configDto)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			configModel := config.Update(db, &config.ConfigModel{DownloadDir: &configDto.DownloadDir})
			c.IndentedJSON(http.StatusOK, ConfigDto{DownloadDir: *configModel.DownloadDir})
		}
	}
}
