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

func (h *GameHandler) TryStartGame(c *gin.Context) {
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

	resp, err := h.gameService.TryStartGame(c.Request.Context(), c.Param("code"), req.VKID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *GameHandler) LoadGameState(c *gin.Context) {
	// TODO: рассчет начальных значений карточек игроков (playerService)

	// обращаю внимание на то, что карточки персонажей должны быть изначально заданы
	// в MVP-версии будет существовать только одна карточка
	resp, err := h.gameService.InitPlayers(c.Request.Context(), c.Param("code"))
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
}

func (h *GameHandler) EndTurn(c *gin.Context) {
	// resp, err := h.gameService.MovePlayer(c.Request.Context(), c.Param("code"))

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
