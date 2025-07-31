package main

import (
	"log"
	"os"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/config"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/handlers"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/migrations"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
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

	err = migrations.RunInitDbMigrations(db)
	if err != nil {
		log.Printf("Error. Failed to migrated db. err: %e", err)
	}

	g := gin.Default()
	handlers.SetupRoutes(g, db)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(g.Run(":" + port))
}
