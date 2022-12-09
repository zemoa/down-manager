package link

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Ref      uuid.UUID `gorm:"type:uuid;index"`
	Link     string    `gorm:"unique;size=1024"`
	Running  bool
	InError  bool
	ErrorMsg *string `gorm:"size=1024"`
}

func Create(link string, db *gorm.DB) *Link {
	uuid, err := uuid.NewUUID()
	if err == nil {
		linkEntity := Link{Link: link, Running: false, InError: false, Ref: uuid}
		db.Create(&linkEntity)
		return &linkEntity
	}
	return nil
}

func GetAll(db *gorm.DB) []Link {
	var links []Link
	db.Find(&links)
	return links
}

func DeleteByRef(ref uuid.UUID, db *gorm.DB) {
	db.Where("ref = ?", ref).Delete(&Link{})
}

func GetByRef(ref uuid.UUID, db *gorm.DB) *Link {
	var link Link
	result := db.Limit(1).Where("ref = ?", ref).Find(&link)
	if result.RowsAffected != 1 {
		return nil
	}
	return &link
}
