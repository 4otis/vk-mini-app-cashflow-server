package dto

type PlayerIsReady struct {
	VKID  int  `json:"vk_id" binding:"required"`
	Ready bool `json:"ready" binding:"required"`
}
