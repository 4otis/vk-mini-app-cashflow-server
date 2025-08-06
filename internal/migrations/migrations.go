package migrations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"gorm.io/gorm"
)

func RunInitDbMigrations(db *gorm.DB) error {
	err := RunAssetMigrations(db)
	if err != nil {
		return err
	}

	err = RunCharaterMigrations(db)
	if err != nil {
		return err
	}

	err = RunIssueMigrations(db)
	if err != nil {
		return err
	}

	err = RunMarketMigrations(db)
	if err != nil {
		return err
	}

	return nil
}

func RunAssetMigrations(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		sqlPath := filepath.Join("internal", "migrations", "sql", "init_assets.sql")
		sqlBytes, err := os.ReadFile(sqlPath)
		if err != nil {
			return err
		}

		if err := tx.Exec(string(sqlBytes)).Error; err != nil {
			return err
		}

		jsonPath := filepath.Join("internal", "migrations", "data", "assets.json")
		jsonBytes, err := os.ReadFile(jsonPath)
		if err != nil {
			return err
		}

		var assets []models.Asset
		if err := json.Unmarshal(jsonBytes, &assets); err != nil {
			return err
		}

		for _, asset := range assets {
			if err := tx.Exec(`
			INSERT INTO assets 
			(created_at, updated_at, title, descr, type_id, price, cashflow)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
				time.Now(), time.Now(),
				asset.Title, asset.Descr,
				asset.TypeID, asset.Price, asset.Cashflow,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func RunCharaterMigrations(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		sqlPath := filepath.Join("internal", "migrations", "sql", "init_characters.sql")
		sqlBytes, err := os.ReadFile(sqlPath)
		if err != nil {
			return err
		}

		if err := tx.Exec(string(sqlBytes)).Error; err != nil {
			return err
		}

		jsonPath := filepath.Join("internal", "migrations", "data", "characters.json")
		jsonBytes, err := os.ReadFile(jsonPath)
		if err != nil {
			return err
		}

		var characters []models.Character
		if err := json.Unmarshal(jsonBytes, &characters); err != nil {
			return err
		}

		for _, c := range characters {
			if err := tx.Exec(`
			INSERT INTO characters 
			(created_at, updated_at, job, salary, taxes, child_expenses, other_expenses)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
				time.Now(), time.Now(),
				c.Job, c.Salary, c.Taxes,
				c.ChildExpenses, c.OtherExpenses,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func RunIssueMigrations(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		sqlPath := filepath.Join("internal", "migrations", "sql", "init_issues.sql")
		sqlBytes, err := os.ReadFile(sqlPath)
		if err != nil {
			return err
		}

		if err := tx.Exec(string(sqlBytes)).Error; err != nil {
			return err
		}

		jsonPath := filepath.Join("internal", "migrations", "data", "issues.json")
		jsonBytes, err := os.ReadFile(jsonPath)
		if err != nil {
			return err
		}

		var issues []models.Issue
		if err := json.Unmarshal(jsonBytes, &issues); err != nil {
			return err
		}

		for _, i := range issues {
			if err := tx.Exec(`
			INSERT INTO issues 
			(created_at, updated_at, title, desc, price)
			VALUES (?, ?, ?, ?, ?)`,
				time.Now(), time.Now(),
				i.Title, i.Descr, i.Price,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func RunMarketMigrations(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		sqlPath := filepath.Join("internal", "migrations", "sql", "init_market.sql")
		sqlBytes, err := os.ReadFile(sqlPath)
		if err != nil {
			return err
		}

		if err := tx.Exec(string(sqlBytes)).Error; err != nil {
			return err
		}

		jsonPath := filepath.Join("internal", "migrations", "data", "market.json")
		jsonBytes, err := os.ReadFile(jsonPath)
		if err != nil {
			return err
		}

		var market []models.Market
		if err := json.Unmarshal(jsonBytes, &market); err != nil {
			return err
		}

		for _, offer := range market {
			if err := tx.Exec(`
			INSERT INTO market 
			(created_at, updated_at, title, desc, type_id, sell_cost)
			VALUES (?, ?, ?, ?, ?, ?)`,
				time.Now(), time.Now(),
				offer.Title, offer.Descr, offer.TypeID, offer.SellCost,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
