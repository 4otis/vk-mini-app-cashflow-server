package repository

import (
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"gorm.io/gorm"
)

type MarketRepository struct {
	db *gorm.DB
}

func NewMarketRepository(db *gorm.DB) *MarketRepository {
	return &MarketRepository{db: db}
}

func (r *MarketRepository) ReadRandom() (market *models.Market, err error) {
	err = r.db.Table("market").Order("RANDOM()").First(&market).Error
	return market, err
}
