package repository

import (
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r SessionRepository) Read(code string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("code = ?", code).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r SessionRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Session{}).Where("id = ?", id).Updates(updates).Error
}

func (r SessionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Session{}, id).Error
}
