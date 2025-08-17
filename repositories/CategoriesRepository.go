package repositories

import (
	"bougette-backend/models"

	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}

// PaginatedCategoriesWithSort returns paginated categories and total count with sorting
func (c *CategoriesRepository) PaginatedCategoriesWithSort(limit, offset int, sort string) ([]models.Categories, int64, error) {
	var categories []models.Categories
	var total int64
	c.db.Model(&models.Categories{}).Count(&total)
	err := c.db.Order(sort).Limit(limit).Offset(offset).Find(&categories).Error
	return categories, total, err
}

func (c *CategoriesRepository) GetCategoryByID(id uint) (*models.Categories, error) {
	var category models.Categories

	err := c.db.First(&category, id).Error
	return &category, err
}

func (c *CategoriesRepository) CreateCategory(category *models.Categories) error {
	return c.db.Create(category).Error
}

func (c *CategoriesRepository) UpdateCategory(category *models.Categories) error {
	if err := c.db.Model(category).Updates(category).Error; err != nil {
		return err
	}

	return nil
}

func (c *CategoriesRepository) DeleteCategory(id uint) error {
	return c.db.Unscoped().Delete(&models.Categories{}, id).Error
}

func (c *CategoriesRepository) FindCategoryByName(name string) (*models.Categories, error) {
	var category models.Categories

	if err := c.db.Where("name = ?", name).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}
