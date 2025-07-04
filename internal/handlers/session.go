package handlers

import (
	"math/rand"
	"net/http"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func generateSessionID() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	seqLen := 16
	out := make([]rune, seqLen)

	for i := range out {
		out[i] = letters[rand.Intn(len(letters))]
	}
	return string(out)
}

func CreateSession(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := generateSessionID()

		session := models.Session{
			ID:       sessionID,
			IsActive: false,
		}

		err := db.Create(&session).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR. Не удалось создать сессию"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"session_id": sessionID,
			"join_link":  "http://vk.com/app123#/join/" + sessionID,
		})

	}
}

func JoinSession(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
