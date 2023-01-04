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

type ConfigRepo struct {
	Db *gorm.DB
}

func (cr *ConfigRepo) Update(config *ConfigModel) ConfigModel {
	var currentConfig Config
	cr.Db.FirstOrCreate(&currentConfig)
	if config.DownloadDir != nil {
		currentConfig.DownloadDir = *config.DownloadDir
	}
	cr.Db.Save(currentConfig)
	return ConfigModel{ID: config.ID, DownloadDir: &currentConfig.DownloadDir}
}

func (cr *ConfigRepo) Get() *Config {
	var config Config
	cr.Db.First(&config)
	return &config
}
