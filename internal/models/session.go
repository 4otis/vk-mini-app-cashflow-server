package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	Code        string   `gorm:"uniqueIndex;size:8"`
	Players     []Player `gorm:"foreignKey:SessionID"`
	IsActive    bool     `gorm:"default:false"`
	CurrentTurn int
	// CreatorID   uint
}
