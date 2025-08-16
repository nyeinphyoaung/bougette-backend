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

func (c *CategoriesService) GetAllCategories() ([]models.Categories, error) {
	return c.CategoriesRepos.GetAllCategories()
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
