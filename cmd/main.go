package main

import (
	"log"
	"os"
	"time"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/config"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/handlers"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/repository"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	db, err := config.InitDB(config.Load().DB)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	// Автомиграции
	if err := db.Migrator().DropTable(
		&models.Session{},
		&models.Player{},
		// Добавьте другие модели...
	); err != nil {
		log.Fatal("Failed to drop all tables:", err)
	}

	if err := db.AutoMigrate(&models.Session{}, &models.Player{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Инициализация зависимостей
	sessionRepo := repository.NewSessionRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	sessionService := services.NewSessionService(sessionRepo, playerRepo)
	sessionHandler := handlers.NewSessionHandler(sessionService)

	// // Создание роутера
	router := gin.Default()
	_ = router

	router.StaticFile("/", "./index.html")           // Для корневого пути
	router.StaticFile("/index.html", "./index.html") // Явно для index.html

	// Маршруты
	router.POST("/sessions", sessionHandler.CreateSession)
	router.POST("/sessions/:code/join", sessionHandler.JoinSession)
	router.GET("/sessions/:code/players", sessionHandler.GetSessionPlayers)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://cshflw.ru/*", "https://vk.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "UPDATE", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(router.Run(":" + port))
}
