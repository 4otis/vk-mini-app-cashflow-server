package repository

import (
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Read(id uint) (asset *models.Asset, err error) {
	err = r.db.First(&asset, id).Error
	return asset, err
}

// func (r *AssetRepository) ReadAll(playerID uint) (assets []*models.Asset, err error) {

// }
