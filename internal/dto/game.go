package dto

type PlayerIsReady struct {
	VKID int `json:"vk_id"`
	// Ready bool `json:"ready" binding:"required"`
}

type PlayersAreReadyRespone struct {
	Ready bool `json:"ready"`
}

type GameStateResponse struct {
	SessionCode string       `json:"code"`
	Players     []PlayerStat `json:"players"`
	CurrentTurn int          `json:"cur_turn"`
	// CurrentCard Card         `json:cur_card`
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

type RollDiceResponse struct {
	Player      PlayerStat `json:"player"`
	CurrentCard Card       `json:"cur_card"`
}

type EndTurnReq struct {
	VKID int `json:"vk_id"`
}

type Card struct {
	Type   string     `json:"type"`
	Asset  AssetCard  `json:"asset,omitempty"`
	Market MarketCard `json:"market,omitempty"`
	Issue  IssueCard  `json:"issue,omitempty"`
}

type CardActionBuyReq struct {
	VKID     int    `json:"vk_id"`
	Title    string `json:"title"`
	Price    int    `json:"price"`
	Cashflow int    `json:"cashflow"`
}

type CardActionSellReq struct {
	VKID     int    `json:"vk_id"`
	Title    string `json:"title"`
	SellCost int    `json:"sell_cost"`
}

type CardActionPayReq struct {
	VKID  int `json:"vk_id"`
	Price int `json:"price"`
}
