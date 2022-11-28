package config

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	DownloadDir string
}

type ConfigModel struct {
	ID          uint
	DownloadDir *string
}

func Update(db *gorm.DB, config *ConfigModel) ConfigModel {
	var currentConfig Config
	db.FirstOrCreate(&currentConfig)
	if config.DownloadDir != nil {
		currentConfig.DownloadDir = *config.DownloadDir
	}
	db.Save(currentConfig)
	return ConfigModel{ID: config.ID, DownloadDir: &currentConfig.DownloadDir}
}

func Get(db *gorm.DB) *Config {
	var config Config
	db.First(&config)
	return &config
}
