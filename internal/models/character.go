package models

import "gorm.io/gorm"

type Character struct {
	gorm.Model
	Job           string `gorm:"size:50;not null" json:"job"`
	Salary        int    `gorm:"type:integer;not null" json:"salary"`
	Taxes         int    `gorm:"type:integer;not null" json:"taxes"`
	ChildExpenses int    `gorm:"type:integer;not null" json:"child_expenses"`
	OtherExpenses int    `gorm:"type:integer;not null" json:"other_expenses"`
}
