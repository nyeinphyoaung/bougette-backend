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
