package dto

type CreateSessionResponse struct {
	Code     string `json:"code"`
	JoinLink string `json:"join_link"`
	// Players  []PlayerResponse `json:"players"`
}
