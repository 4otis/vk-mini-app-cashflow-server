package repository

import (
	"fmt"
	"math/rand/v2"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/dto"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"gorm.io/gorm"
)

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

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

func (r PlayerRepository) ReadByVKID(vkID int) (player *models.Player, err error) {
	err = r.db.Where("vk_id = ?", vkID).First(&player).Error
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

func (r PlayerRepository) InitPlayer(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Verify player exists
		var player models.Player
		if err := tx.First(&player, id).Error; err != nil {
			return fmt.Errorf("player not found: %w", err)
		}

		var character models.Character
		if err := tx.Order("RANDOM()").First(&character).Error; err != nil {
			return fmt.Errorf("failed to get random character: %w", err)
		}

		totalIncome := 0 + character.Salary
		totalExpenses := character.Taxes + character.ChildExpenses*player.ChildAmount + character.OtherExpenses

		// 5. Prepare player initialization params
		params := map[string]interface{}{
			"character_id":   character.ID,
			"passive_income": 0,
			"total_income":   totalIncome,
			"total_expenses": totalExpenses,
			"cashflow":       totalIncome - totalExpenses,
			"balance":        0,
			"bank_loan":      0,
		}

		// 6. Update player fields
		if err := tx.Model(&models.Player{}).Where("id = ?", id).Updates(params).Error; err != nil {
			return fmt.Errorf("failed to update player: %w", err)
		}

		return nil
	})
}

func isPayday(pos int, value int) bool {
	return value >= (8 - pos%8)
}

func paydayPlayer(tx *gorm.DB, player *models.Player) error {
	newBalance := player.Balance + player.Cashflow

	return tx.Model(player).Update("balance", newBalance).Error
}

func (r PlayerRepository) MovePlayer(VKID int, value int) (*models.Player, error) {
	var resultPlayer *models.Player

	err := r.db.Transaction(func(tx *gorm.DB) error {
		player := &models.Player{}
		if err := tx.Where("vk_id = ?", VKID).First(player).Error; err != nil {
			return fmt.Errorf("failed to find player: %w", err)
		}

		if isPayday(player.Position, value) {
			if err := paydayPlayer(tx, player); err != nil {
				return fmt.Errorf("payday failed: %w", err)
			}
		}

		newPosition := (player.Position + value) % 24
		if err := tx.Model(player).Update("position", newPosition).Error; err != nil {
			return fmt.Errorf("failed to update position: %w", err)
		}

		resultPlayer = player
		return nil
	})

	return resultPlayer, err
}

func (r PlayerRepository) BuyAsset(req *dto.CardActionBuyReq) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		asset := &models.Asset{}
		err := tx.Model(asset).Where("title = ? AND price = ? AND cashflow = ?",
			req.Title, req.Price, req.Cashflow).First(asset).Error
		if err != nil {
			return fmt.Errorf("error while reading assetID by args: %w", err)
		}

		player := &models.Player{}
		err = tx.Model(player).Where("vk_id = ?", req.VKID).First(&player).Error
		if err != nil {
			return fmt.Errorf("error while reading player by VKID: %w", err)
		}

		if player.Balance < req.Price {
			return fmt.Errorf("error player.Balance is not enough")
		} else {
			player.Balance -= req.Price
			player.PassiveIncome += req.Cashflow
			player.TotalIncome += req.Cashflow
			player.Cashflow = player.TotalIncome - player.TotalExpenses
		}

		err = tx.Model(player).Where("id = ?", player.ID).Updates(player).Error
		if err != nil {
			return fmt.Errorf("error while updating player: %w", err)
		}

		err = tx.Exec(
			"INSERT INTO players_assets (asset_id, player_id) VALUES (?, ?)",
			asset.ID,
			player.ID,
		).Error
		if err != nil {
			return fmt.Errorf("error while inserting into players_assets: %w", err)
		}

		return nil
	})
}

func (r PlayerRepository) SellAsset(req *dto.CardActionSellReq) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		player := &models.Player{}
		if err := tx.Where("vk_id = ?", req.VKID).First(player).Error; err != nil {
			return fmt.Errorf("error finding player: %w", err)
		}

		var assetsToSell []models.Asset
		for _, assetStat := range req.Assets {
			var asset models.Asset
			if err := tx.Where("title = ? AND price = ? AND cashflow = ?",
				assetStat.Title, assetStat.Price, assetStat.Cashflow).First(&asset).Error; err != nil {
				return fmt.Errorf("error finding asset %s: %w", assetStat.Title, err)
			}

			var count int64
			if err := tx.Table("players_assets").
				Where("player_id = ? AND asset_id = ?", player.ID, asset.ID).
				Count(&count).Error; err != nil || count == 0 {
				return fmt.Errorf("player doesn't own asset %s or error checking ownership: %w", assetStat.Title, err)
			}

			assetsToSell = append(assetsToSell, asset)
		}

		totalCashflowReduction := 0
		for _, asset := range assetsToSell {
			totalCashflowReduction += asset.Cashflow
		}

		totalIncome := len(req.Assets) * req.SellCost
		player.Balance += totalIncome
		player.PassiveIncome -= totalCashflowReduction
		player.TotalIncome -= totalCashflowReduction
		player.Cashflow = player.TotalIncome - player.TotalExpenses

		err := tx.Model(player).Where("id = ?", player.ID).Updates(player).Error
		if err != nil {
			return fmt.Errorf("error while updating player: %w", err)
		}

		for _, asset := range assetsToSell {
			if err := tx.Exec(
				"DELETE FROM players_assets WHERE player_id = ? AND asset_id = ?",
				player.ID,
				asset.ID,
			).Error; err != nil {
				return fmt.Errorf("error deleting player_asset relation: %w", err)
			}
		}

		return nil
	})
}
