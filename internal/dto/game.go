package dto

type PlayerIsReady struct {
	VKID int `json:"vk_id"`
	// Ready bool `json:"ready" binding:"required"`
}

type GameStateResponse struct {
	SessionCode string       `json:"code"`
	Players     []PlayerStat `json:"players"`
	CurrentTurn int          `json:"cur_turn"`
	// Phase       GamePhase    `json:"phase"` // "waiting", "rolling", "trading", "end_turn"
	// Board       []Cell       `json:"board,omitempty"`
	// Deck        struct {
	// 	Opportunities []Card `json:"opportunities"`
	// 	Market        []Card `json:"market"`
	// } `json:"deck,omitempty"`
	// Log []GameLogEntry `json:"log"`
}

type RollDiceReq struct {
	VKID      int `json:"vk_id"`
	DiceValue int `json:"dice_value"`
}

type EndTurnReq struct {
	VKID int `json:"vk_id"`
}
