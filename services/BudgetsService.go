package services

import (
	"bougette-backend/models"
	"bougette-backend/repositories"
)

type BudgetsService struct {
	BudgetsRepos *repositories.BudgetsRepository
}

func NewBudgetsService(budgetsRepo *repositories.BudgetsRepository) *BudgetsService {
	return &BudgetsService{BudgetsRepos: budgetsRepo}
}

func (b *BudgetsService) CreateBudgets(budget *models.Budgets) error {
	return b.BudgetsRepos.CreateBudgets(budget)
}

func (b *BudgetsService) CheckBudgetsExit(UserID uint, month uint, year uint16, slug string) (bool, error) {
	budget, err := b.BudgetsRepos.CheckBudgetsExit(UserID, month, year, slug)
	if err != nil {
		return false, err
	}
	return budget != nil, nil
}

func (b *BudgetsService) UpdateBudgetCategories(budgetID uint, categories []models.Categories) error {
	return b.BudgetsRepos.UpdateBudgetCategories(budgetID, categories)
}

func (b *BudgetsService) GetBudgetWithCategories(budgetID uint) (*models.Budgets, error) {
	return b.BudgetsRepos.GetBudgetWithCategories(budgetID)
}

func (b *BudgetsService) GetPaginatedBudgetsByUserID(userID uint, limit, offset int, sort string) ([]models.Budgets, int64, error) {
	return b.BudgetsRepos.GetPaginatedBudgetsByUserID(userID, limit, offset, sort)
}

func (b *BudgetsService) UpdateBudget(budget *models.Budgets) error {
	return b.BudgetsRepos.UpdateBudget(budget)
}

func (b *BudgetsService) CheckBudgetsExitExcludingID(UserID uint, month uint, year uint16, slug string, excludeID uint) (bool, error) {
	budget, err := b.BudgetsRepos.CheckBudgetsExitExcludingID(UserID, month, year, slug, excludeID)
	if err != nil {
		return false, err
	}
	return budget != nil, nil
}

func (b *BudgetsService) GetBudgetByID(id uint) (*models.Budgets, error) {
	return b.BudgetsRepos.GetBudgetByID(id)
}

func (b *BudgetsService) DeleteBudget(id uint) error {
	return b.BudgetsRepos.DeleteBudget(id)
}

func (b *BudgetsService) CreateBudgetWithCategories(budget *models.Budgets, categoryIDs []uint64) error {
	return b.BudgetsRepos.CreateBudgetWithCategories(budget, categoryIDs)
}

func (b *BudgetsService) UpdateBudgetWithCategories(budget *models.Budgets, categoryIDs []uint64) error {
	return b.BudgetsRepos.UpdateBudgetWithCategories(budget, categoryIDs)
}
