package handlers

import (
	"time"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/repository"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(g *gin.Engine, db *gorm.DB) {
	// Инициализация зависимостей
	sessionRepo := repository.NewSessionRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	assetRepo := repository.NewAssetRepository(db)

	sessionService := services.NewSessionService(sessionRepo, playerRepo)
	gameService := services.NewGameService(sessionRepo, playerRepo, assetRepo)

	sessionHandler := NewSessionHandler(sessionService)
	gameHandler := NewGameHandler(gameService)

	g.StaticFile("/", "./index.html")           // Для корневого пути
	g.StaticFile("/index.html", "./index.html") // Явно для index.html

	// Маршруты
	g.POST("/sessions", sessionHandler.CreateSession)
	g.POST("/sessions/:code/join", sessionHandler.JoinSession)
	g.GET("/sessions/:code/players", sessionHandler.GetSessionPlayers)
	g.PATCH("/game/:code/ready", gameHandler.TryStartGame)
	g.GET("/game/:code/state", gameHandler.LoadGameState)

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://cshflw.ru/*", "https://vk.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
