package link

import "gorm.io/gorm"

type Link struct {
	gorm.Model
	Link     string `gorm:"unique;size=1024;index"`
	Running  bool
	InError  bool
	ErrorMsg *string `gorm:"size=1024"`
}

func Create(link string, db *gorm.DB) *Link {
	linkEntity := Link{Link: link, Running: false, InError: false}
	db.Create(&linkEntity)
	return &linkEntity
}

func GetAll(db *gorm.DB) []Link {
	var links []Link
	db.Find(&links)
	return links
}
