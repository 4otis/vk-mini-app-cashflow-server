package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	VKID          int      `gorm:"uniqueIndex:idx_vkid_session;not null" json:"vk_id"`
	SessionID     uint     `gorm:"index;not null" json:"session_id"`
	Nickname      string   `gorm:"size:50;not null" json:"nickname"`
	Ready         bool     `gorm:"default:false" json:"ready"`
	Position      int      `gorm:"default:0" json:"position"`
	CharacterID   uint     `gorm:"index" json:"character_id"`
	PassiveIncome float64  `gorm:"type:decimal(10,2);default:0" json:"passive_income"`
	TotalIncome   float64  `gorm:"type:decimal(10,2);default:0" json:"total_income"`
	TotalExpenses float64  `gorm:"type:decimal(10,2);default:0" json:"total_expenses"`
	CashFlow      float64  `gorm:"type:decimal(10,2);default:0" json:"cashflow"`
	Balance       float64  `gorm:"type:decimal(10,2);default:0" json:"balance"`
	BankLoan      float64  `gorm:"type:decimal(10,2);default:0" json:"bank_loan"`
	Asset         []*Asset `gorm:"many2many:players_assets;"`
}
