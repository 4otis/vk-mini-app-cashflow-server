package repository

import (
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"gorm.io/gorm"
)

type IssueRepository struct {
	db *gorm.DB
}

func NewIssueRepository(db *gorm.DB) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) ReadRandom() (issue *models.Issue, err error) {
	err = r.db.Order("RANDOM()").First(&issue).Error
	return issue, err
}
