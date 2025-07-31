package dto

type CreatePlayerRequest struct {
	VKID     int    `json:"vk_id" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type PlayerResponse struct {
	ID            uint   `json:"id"`
	VKID          int    `json:"vk_id"`
	Nickname      string `json:"nickname"`
	Ready         bool   `json:"ready"`
	Position      int    `json:"position"`
	CharacterID   uint   `gorm:"index" json:"character_id"`
	PassiveIncome int    `json:"passive_income"`
	TotalIncome   int    `json:"total_income"`
	Cashflow      int    `json:"cashflow"`
	Balance       int    `json:"balance"`
	BankLoan      int    `json:"bank_loan"`
}

type PlayerStat struct {
	VKID          int         `json:"vk_id"`
	Nickname      string      `json:"nickname"`
	PassiveIncome int         `json:"passive_income"`
	TotalIncome   int         `json:"total_income"`
	TotalExpenses int         `json:"total_expenses"`
	Cashflow      int         `json:"cashflow"`
	Position      int         `json:"position"`
	Assets        []AssetStat `json:"assets"`
	Balance       int         `json:"balance"`
	ChildAmount   int         `json:"child_amount"`
	BankLoan      int         `json:"bank_loan"`
	// Liabilities   []Liability `json:"liabilities"`
	// IsBankrupt    bool        `json:"is_bankrupt"`
}
