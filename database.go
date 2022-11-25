package database

import "gorm.io/gorm"

type Link struct {
	gorm.Model
	link string
}

func Init() {

}
