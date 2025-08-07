package models

import "gorm.io/gorm"

type Asset struct {
	gorm.Model
	Title    string    `gorm:"size:50;not null" json:"title"`
	Descr    string    `gorm:"size:100;not null" json:"descr"`
	TypeID   int       `gorm:"index;not null" json:"type_id"`
	Price    int       `gorm:"type:integer;default:0" json:"price"`
	Cashflow int       `gorm:"type:integer;default:0" json:"cashflow"`
	Players  []*Player `gorm:"many2many:players_assets;OnDelete:CASCADE"`
}
