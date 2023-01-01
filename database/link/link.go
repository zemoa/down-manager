package link

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Ref            uuid.UUID `gorm:"type:uuid;index"`
	Link           string    `gorm:"unique;size=1024"`
	Filename       string    `gorm:"size=128"`
	Rangesupported bool
	Size           uint64
	Downloaded     uint64
	Running        bool
	InError        bool
	ErrorMsg       *string `gorm:"size=1024"`
}

type LinkRepo struct {
	Db *gorm.DB
}

func (lr *LinkRepo) Create(link string) *Link {
	uuid, err := uuid.NewUUID()
	if err == nil {
		linkEntity := Link{Link: link, Running: false, InError: false, Ref: uuid}
		lr.Db.Create(&linkEntity)
		return &linkEntity
	}
	return nil
}

func (lr *LinkRepo) Update(link *Link) {
	lr.Db.Save(link)
}

func (lr *LinkRepo) GetAll() []Link {
	var links []Link
	lr.Db.Find(&links)
	return links
}

func (lr *LinkRepo) DeleteByRef(ref uuid.UUID) {
	lr.Db.Where("ref = ?", ref).Delete(&Link{})
}

func (lr *LinkRepo) GetByRef(ref uuid.UUID) *Link {
	var link Link
	result := lr.Db.Limit(1).Where("ref = ?", ref).Find(&link)
	if result.RowsAffected != 1 {
		return nil
	}
	return &link
}

func (lr *LinkRepo) UpdateDownloaded(ref uuid.UUID, downloaded uint64) {
	lr.Db.Model(&Link{}).Where("ref = ?", ref).Update("downloaded", downloaded)
}
