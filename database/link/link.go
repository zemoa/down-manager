package link

import "gorm.io/gorm"

type Link struct {
	gorm.Model
	link string `gorm:"unique;size=1024;index"`
}

func Create(link string, db *gorm.DB) *Link {
	linkEntity := Link{link: link}
	db.Create(&linkEntity)
	return &linkEntity
}

func GetAll(db *gorm.DB) []Link {
	var links []Link
	db.Find(&links)
	return links
}
