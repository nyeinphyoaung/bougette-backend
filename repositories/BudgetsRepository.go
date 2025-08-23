package repositories

import (
	"bougette-backend/models"

	"gorm.io/gorm"
)

type BudgetsRepository struct {
	db *gorm.DB
}

func NewBudgetsRepository(db *gorm.DB) *BudgetsRepository {
	return &BudgetsRepository{db: db}
}

func (b *BudgetsRepository) CreateBudgets(budget *models.Budgets) error {
	return b.db.Create(budget).Error
}

func (b *BudgetsRepository) CheckBudgetsExit(UserID uint, month uint, year uint16, slug string) (*models.Budgets, error) {
	var budget models.Budgets

	if err := b.db.Where("user_id = ? AND month = ? AND year = ? AND slug = ?", UserID, month, year, slug).First(&budget).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &budget, nil
}

func (b *BudgetsRepository) UpdateBudgetCategories(budgetID uint, categories []models.Categories) error {
	var budget models.Budgets
	if err := b.db.First(&budget, budgetID).Error; err != nil {
		return err
	}

	return b.db.Model(&budget).Association("Categories").Replace(categories)
}

func (b *BudgetsRepository) GetBudgetWithCategories(budgetID uint) (*models.Budgets, error) {
	var budget models.Budgets
	err := b.db.Preload("Categories").First(&budget, budgetID).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}
