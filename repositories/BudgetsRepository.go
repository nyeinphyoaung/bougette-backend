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

func (b *BudgetsRepository) GetPaginatedBudgetsByUserID(userID uint, limit, offset int, sort string) ([]models.Budgets, int64, error) {
	var budgets []models.Budgets
	var total int64

	// Count total budgets for this user
	err := b.db.Model(&models.Budgets{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated budgets with categories preloaded
	err = b.db.Where("user_id = ?", userID).Preload("Categories").Order(sort).Limit(limit).Offset(offset).Find(&budgets).Error
	if err != nil {
		return nil, 0, err
	}

	return budgets, total, nil
}

func (b *BudgetsRepository) UpdateBudget(budget *models.Budgets) error {
	if err := b.db.Model(budget).Updates(budget).Error; err != nil {
		return err
	}

	return nil
}

// CheckBudgetsExitExcludingID checks for an existing budget matching the
// provided user/month/year/slug combination, while excluding the given budget ID.
func (b *BudgetsRepository) CheckBudgetsExitExcludingID(UserID uint, month uint, year uint16, slug string, excludeID uint) (*models.Budgets, error) {
	var budget models.Budgets

	if err := b.db.Where("user_id = ? AND month = ? AND year = ? AND slug = ? AND id <> ?", UserID, month, year, slug, excludeID).First(&budget).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &budget, nil
}

func (b *BudgetsRepository) GetBudgetByID(id uint) (*models.Budgets, error) {
	var budget models.Budgets

	err := b.db.First(&budget, id).Error
	return &budget, err
}
