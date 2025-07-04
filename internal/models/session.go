package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	ID          string `gorm:"uniqueIndex"`
	Players     []Player
	IsActive    bool
	CurrentTurn uint
}

type Player struct {
	gorm.Model
	SessionID string
	VKID      uint `gorm:"index"`
	Nickname  string
	Ready     bool
	Balance   int `gorm:"default:1000"`
	Position  int `gorm:"default:0"`
}
