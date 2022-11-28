package service

import (
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
