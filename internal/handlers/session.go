package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/dto"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/services"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService *services.SessionService
}

func NewSessionHandler(sessionService *services.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req dto.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR. BadRequest: (VKID=%d; Nickname=%s)\n", req.VKID, req.Nickname)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.VKID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	resp, err := h.sessionService.CreateSession(c.Request.Context(), req.VKID, req.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *SessionHandler) JoinSession(c *gin.Context) {
	code := c.Param("code")

	var req dto.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.VKID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	player, err := h.sessionService.JoinSession(c.Request.Context(), code, req.VKID, req.Nickname)
	if err != nil {
		switch err {
		case services.ErrSessionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		case services.ErrInvalidCode:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session code"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, player)
}

func (h *SessionHandler) GetSessionPlayers(c *gin.Context) {
	code := c.Param("code")
	// if len(code) != 6 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session code format"})
	// 	return
	// }

	players, err := h.sessionService.GetSessionPlayers(c.Request.Context(), code)
	if err != nil {
		switch err {
		case services.ErrSessionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, players)
}

func (h *SessionHandler) DeleteSession(c *gin.Context) {
	code := c.Param("code")

	// Удаляем сессию через сервис
	err := h.sessionService.DeleteSession(c.Request.Context(), code)
	if err != nil {
		// Определяем тип ошибки для соответствующего HTTP-статуса
		switch {
		case errors.Is(err, services.ErrSessionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Сессия не найдена"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить сессию"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Сессия успешно удалена",
		"code":    code,
	})
}
