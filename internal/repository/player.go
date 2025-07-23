package repository

import (
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"gorm.io/gorm"
)

type PlayerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r PlayerRepository) Create(player *models.Player) error {
	return r.db.Create(player).Error
}

func (r PlayerRepository) Read(id uint) (player *models.Player, err error) {
	err = r.db.First(&player, id).Error
	return player, err
}

func (r PlayerRepository) ReadAll(sessionID uint) ([]models.Player, error) {
	var players []models.Player
	err := r.db.Where("session_id = ?", sessionID).Find(&players).Error
	return players, err
}

func (r PlayerRepository) Update(id uint, newPlayer *models.Player) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var player models.Player
		if err := tx.First(&player, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&player).Updates(newPlayer).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r PlayerRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Player{}).Where("id = ?", id).Updates(updates).Error
}

func (r PlayerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Player{}, id).Error
}
