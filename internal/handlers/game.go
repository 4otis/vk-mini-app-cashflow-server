package handlers

import (
	"log"
	"net/http"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/dto"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/services"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (h *GameHandler) PlayerIsReady(c *gin.Context) {
	var req dto.PlayerIsReady
	if err := c.ShouldBindJSON(&req); err != nil {
		// log.Printf("ERROR. BadRequest: (VKID=%d;)\n", req.VKID)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.VKID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	resp, err := h.gameService.PlayerIsReady(c.Request.Context(), c.Param("code"), req.VKID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *GameHandler) PlayersAreReady(c *gin.Context) {
	resp, err := h.gameService.ArePlayersReady(c.Request.Context(), c.Param("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *GameHandler) InitGameState(c *gin.Context) {
	resp, err := h.gameService.InitGameState(c.Request.Context(), c.Param("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *GameHandler) LoadGameState(c *gin.Context) {
	resp, err := h.gameService.LoadGameState(c.Request.Context(), c.Param("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *GameHandler) RollDice(c *gin.Context) {
	// resp, err := h.gameService.MovePlayer(c.Request.Context(), c.Param("code"))

	var req dto.RollDiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// log.Printf("ERROR. BadRequest: (VKID=%d;)\n", req.VKID)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Player: %d, value: %d;\n", req.VKID, req.DiceValue)

	resp, err := h.gameService.RollDice(c.Request.Context(), c.Param("code"), req.VKID, req.DiceValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *GameHandler) EndTurn(c *gin.Context) {
	var req dto.EndTurnReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// log.Printf("ERROR. BadRequest: (VKID=%d;)\n", req.VKID)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// log.Printf("Player: %d, value: %d;\n", req.VKID, req.DiceValue)

	err := h.gameService.EndTurn(c.Request.Context(), c.Param("code"), req.VKID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *GameHandler) CardActionBuy(c *gin.Context) {
	log.Printf("CardActionBuy was TRIGGERED.")

	var req dto.CardActionBuyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.gameService.BuyAsset(c.Request.Context(), c.Param("code"), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *GameHandler) CardActionSell(c *gin.Context) {
	log.Printf("CardActionSell was TRIGGERED.")

	var req dto.CardActionSellReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.gameService.SellAsset(c.Request.Context(), c.Param("code"), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *GameHandler) CardActionPay(c *gin.Context) {
	log.Printf("CardActionPay was TRIGGERED.")
}

func (h *GameHandler) CardActionChild(c *gin.Context) {
	log.Printf("CardActionChild was TRIGGERED.")
}
