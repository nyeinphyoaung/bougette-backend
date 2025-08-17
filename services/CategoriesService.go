package services

import (
	"bougette-backend/models"
	"bougette-backend/repositories"
)

type CategoriesService struct {
	CategoriesRepos *repositories.CategoriesRepository
}

func NewCategoriesService(categoriesRepo *repositories.CategoriesRepository) *CategoriesService {
	return &CategoriesService{CategoriesRepos: categoriesRepo}
}

// PaginatedCategoriesWithSort returns paginated categories and total count with sorting
func (c *CategoriesService) PaginatedCategoriesWithSort(limit, offset int, sort string) ([]models.Categories, int64, error) {
	return c.CategoriesRepos.PaginatedCategoriesWithSort(limit, offset, sort)
}

func (c *CategoriesService) GetCategoryByID(id uint) (*models.Categories, error) {
	return c.CategoriesRepos.GetCategoryByID(id)
}

func (c *CategoriesService) CreateCategory(category *models.Categories) error {
	return c.CategoriesRepos.CreateCategory(category)
}

func (c *CategoriesService) UpdateCategory(category *models.Categories) error {
	return c.CategoriesRepos.UpdateCategory(category)
}

func (c *CategoriesService) DeleteCategory(id uint) error {
	return c.CategoriesRepos.DeleteCategory(id)
}

func (u *CategoriesService) CheckCategoryExits(name string) (bool, error) {
	category, err := u.CategoriesRepos.FindCategoryByName(name)
	if err != nil {
		return false, err
	}
	return category != nil, nil
}

func (c *CategoriesService) FindCategoryByName(name string) (*models.Categories, error) {
	return c.CategoriesRepos.FindCategoryByName(name)
}
