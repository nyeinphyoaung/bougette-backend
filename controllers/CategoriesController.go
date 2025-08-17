package controllers

import (
	"bougette-backend/common"
	"bougette-backend/dtos"
	"bougette-backend/models"
	"bougette-backend/services"
	"bougette-backend/validation"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type CategoriesController struct {
	CategoriesService *services.CategoriesService
}

func NewCategoriesController(categoriesService *services.CategoriesService) *CategoriesController {
	return &CategoriesController{CategoriesService: categoriesService}
}

// Efficient pagination implemented here.
// For more advanced options, refer to the official GORM Scopes pagination documentation.
func (c *CategoriesController) GetPaginatedCategories(ctx echo.Context) error {
	// page = 2
	// limit = 5
	// offset = (2 - 1) * 5 = 5 → skip first 5 items
	page := 1
	limit := 10
	sort := "id desc" // name desc dynamically like that
	if p := ctx.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := ctx.QueryParam("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if s := ctx.QueryParam("sort"); s != "" {
		sort = s
	}
	offset := (page - 1) * limit
	categories, total, err := c.CategoriesService.PaginatedCategoriesWithSort(limit, offset, sort)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, err.Error())
	}
	response := map[string]interface{}{
		"categories": categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"sort":       sort,
	}
	return common.SendSuccessResponse(ctx, "Categories retrieved successfully", response)
}

func (c *CategoriesController) GetCategoryByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid category ID")
	}

	category, err := c.CategoriesService.GetCategoryByID(uint(id))
	if err != nil {
		return common.SendNotFoundResponse(ctx, "Category not found")
	}

	return common.SendSuccessResponse(ctx, "Category retrieved successfully", category)
}

func (c *CategoriesController) CreateCategory(ctx echo.Context) error {
	request := new(dtos.CreateCategoryRequestDTO)
	if err := ctx.Bind(request); err != nil {
		return common.SendBadRequestResponse(ctx, err.Error())
	}

	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		return common.SendFailedValidationResponse(ctx, validationErrors)
	}

	exits, err := c.CategoriesService.CheckCategoryExits(request.Name)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, err.Error())
	}

	if exits {
		return common.SendBadRequestResponse(ctx, "Category with this name already exits")
	}

	slug := strings.ToLower(request.Name)
	slug = strings.Replace(slug, " ", "_", -1)

	category := models.Categories{
		Name: request.Name,
		Slug: slug,
	}
	if err := c.CategoriesService.CreateCategory(&category); err != nil {
		return common.SendNotFoundResponse(ctx, err.Error())
	}

	return common.SendSuccessResponse(ctx, "Category created successfully", category)
}

func (c *CategoriesController) UpdateCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid category ID")
	}

	exitstingCategory, err := c.CategoriesService.GetCategoryByID(uint(id))
	if err != nil {
		return common.SendNotFoundResponse(ctx, "Category not found")
	}

	updateCategory := new(dtos.UpdateCategoryRequestDTO)
	if err := ctx.Bind(updateCategory); err != nil {
		return common.SendBadRequestResponse(ctx, err.Error())
	}

	if updateCategory.Name != nil {
		exitstingCategory.Name = *updateCategory.Name
	}
	if updateCategory.Slug != nil {
		exitstingCategory.Slug = *updateCategory.Slug
	}
	if updateCategory.IsCustom != nil {
		exitstingCategory.IsCustom = *updateCategory.IsCustom
	}

	if err := c.CategoriesService.UpdateCategory(exitstingCategory); err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to update category")
	}

	return common.SendSuccessResponse(ctx, "Category has been updated successfully", exitstingCategory)
}

func (c *CategoriesController) DeleteCategory(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid category ID")
	}

	err = c.CategoriesService.DeleteCategory(uint(id))
	if err != nil {
		return common.SendNotFoundResponse(ctx, "Category not found")
	}

	return common.SendSuccessResponse(ctx, "Category deleted successfully", nil)
}
