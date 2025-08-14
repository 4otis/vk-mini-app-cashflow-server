package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	VKID          int    `gorm:"uniqueIndex:idx_vkid_session,where:deleted_at IS NULL;not null" json:"vk_id"`
	SessionID     uint   `gorm:"index;not null" json:"session_id"`
	Nickname      string `gorm:"size:50;not null" json:"nickname"`
	Ready         bool   `gorm:"default:false" json:"ready"`
	Position      int    `gorm:"default:0" json:"position"`
	CharacterID   uint   `gorm:"index" json:"character_id"`
	PassiveIncome int    `gorm:"type:integer;default:0" json:"passive_income"`
	TotalIncome   int    `gorm:"type:integer;default:0" json:"total_income"`
	TotalExpenses int    `gorm:"type:integer;default:0" json:"total_expenses"`
	Cashflow      int    `gorm:"type:integer;default:0" json:"cashflow"`
	Balance       int    `gorm:"type:integer;default:0" json:"balance"`
	BankLoan      int    `gorm:"type:integer;default:0" json:"bank_loan"`
	ChildAmount   int    `gorm:"type:integer;default:0" json:"child_amount"`
	// Asset         []*Asset  `gorm:"many2many:players_assets;OnDelete:CASCADE"`
	Character Character `gorm:"foreignKey:CharacterID;constraint:OnDelete:CASCADE" json:"character"`
	IsPayday  bool
}
