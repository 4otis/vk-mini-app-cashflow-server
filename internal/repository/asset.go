package repository

import (
	"fmt"

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

func (r *AssetRepository) ReadRandom() (asset *models.Asset, err error) {
	err = r.db.Order("RANDOM()").First(&asset).Error
	return asset, err
}

func (r *AssetRepository) ReadAllByPlayerID(id uint) ([]*models.Asset, error) {
	var assets []*models.Asset
	err := r.db.
		Joins("JOIN players_assets ON players_assets.asset_id = assets.id").
		Where("players_assets.player_id = ? AND assets.deleted_at IS NULL", id).
		Find(&assets).Error
	if err != nil {
		return assets, fmt.Errorf("error: failed to load assets: %v", err)
	}
	return assets, nil
}
