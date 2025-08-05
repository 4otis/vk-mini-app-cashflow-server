package models

import "gorm.io/gorm"

type Issue struct {
	gorm.Model
	Title string `gorm:"size:50;not null" json:"title"`
	Descr string `gorm:"size:100;not null" json:"descr"`
	Price int    `gorm:"type:integer;default:0" json:"price"`
}
