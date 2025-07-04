package main

import (
	"log"
	"os"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// Инициализация БД
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автомиграции
	db.AutoMigrate(&models.Session{}, &models.Player{})

	// Создание сервера Gin
	r := gin.Default()

	// Регистрация роутов
	registerAPIRoutes(r, db)

	// Запуск сервера
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func registerAPIRoutes(r *gin.Engine, db *gorm.DB) {
	// Группа роутов для сессий
	sessionGroup := r.Group("/session")
	{
		sessionGroup.POST("/create", handlers.CreateSession(db))
		sessionGroup.POST("/join", handlers.JoinSession(db))
	}

	// Группа роутов для игры
	gameGroup := r.Group("/cashflow/:id")
	{
		gameGroup.GET("/players", handlers.GetPlayers(db))
		gameGroup.POST("/ready", handlers.ToggleReady(db))
		gameGroup.POST("/move", handlers.MakeMove(db))
	}
}
