package dto

type CreatePlayerRequest struct {
	VKID     int    `json:"vk_id" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type PlayerResponse struct {
	ID            uint    `json:"id"`
	VKID          int     `json:"vk_id"`
	Nickname      string  `json:"nickname"`
	Ready         bool    `json:"ready"`
	Position      int     `json:"position"`
	CharacterID   uint    `gorm:"index" json:"character_id"`
	PassiveIncome float64 `json:"passive_income"`
	TotalIncome   float64 `json:"total_income"`
	CashFlow      float64 `json:"cash_flow"`
	Balance       float64 `json:"balance"`
	BankLoan      float64 `json:"bank_loan"`
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
	// Liabilities   []Liability `json:"liabilities"`
	// IsBankrupt    bool        `json:"is_bankrupt"`
	// Children      int         `json:"children"`
}
