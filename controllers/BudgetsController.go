package controllers

import (
	"bougette-backend/common"
	"bougette-backend/dtos"
	"bougette-backend/models"
	"bougette-backend/services"
	"bougette-backend/validation"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type BudgetsController struct {
	BudgetsService    *services.BudgetsService
	CategoriesService *services.CategoriesService
}

func NewBudgetsController(budgetsService *services.BudgetsService, categoriesService *services.CategoriesService) *BudgetsController {
	return &BudgetsController{BudgetsService: budgetsService, CategoriesService: categoriesService}
}

func (b *BudgetsController) CreateBudgets(ctx echo.Context) error {
	userID, ok := ctx.Get("user").(uint)
	// fmt.Println("User ID:", userID)
	if !ok {
		return common.SendInternalServerErrorResponse(ctx, "User authentication required")
	}

	request := new(dtos.CreateBudgetRequestDTO)
	if err := ctx.Bind(request); err != nil {
		return common.SendBadRequestResponse(ctx, err.Error())
	}

	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		return common.SendFailedValidationResponse(ctx, validationErrors)
	}

	slug := strings.ToLower(request.Title)
	slug = strings.Replace(slug, " ", "_", -1)

	budget := models.Budgets{
		Title:       request.Title,
		Slug:        slug,
		UserID:      userID,
		Amount:      request.Amount,
		Description: request.Description,
	}
	if request.Date == "" {
		currentDate := time.Now()
		budget.Date = currentDate
	} else {
		// Parse the date string - support both YYYY-MM-DD and YYYY-MM-DD HH:MM:SS formats
		var parsedDate time.Time
		var err error

		if len(request.Date) == 10 { // YYYY-MM-DD format
			parsedDate, err = time.Parse("2006-01-02", request.Date)
		} else { // YYYY-MM-DD HH:MM:SS format
			parsedDate, err = time.Parse("2006-01-02 15:04:05", request.Date)
		}

		if err != nil {
			return common.SendBadRequestResponse(ctx, "Invalid date format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS")
		}
		budget.Date = parsedDate
	}

	month := uint(budget.Date.Month())
	year := uint16(budget.Date.Year())

	budget.Month = month
	budget.Year = year

	// Debug logging
	fmt.Printf("Debug - Date: %v, Month: %d, Year: %d, Slug: %s\n", budget.Date, month, year, budget.Slug)

	exits, err := b.BudgetsService.CheckBudgetsExit(budget.UserID, budget.Month, budget.Year, budget.Slug)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, err.Error())
	}

	if exits {
		return common.SendBadRequestResponse(ctx, "Budget with this userId, month, year and slug already exits")
	}

	if err := b.BudgetsService.CreateBudgets(&budget); err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Budget could not be created")
	}

	// get all categories and attach categories association
	if len(request.Categories) > 0 {
		for _, catID := range request.Categories {
			if _, err := b.CategoriesService.GetCategoryByID(uint(catID)); err != nil {
				return common.SendBadRequestResponse(ctx, fmt.Sprintf("Category with ID %d does not exist", catID))
			}
		}
		categories, err := b.CategoriesService.GetCategoriesByIDs(request.Categories)
		if err != nil {
			return common.SendInternalServerErrorResponse(ctx, "Failed to fetch categories")
		}

		// Update the budget with categories association
		if err := b.BudgetsService.UpdateBudgetCategories(budget.ID, categories); err != nil {
			return common.SendInternalServerErrorResponse(ctx, "Failed to attach categories to budget")
		}
	}

	// Fetch the budget with categories loaded for response
	createdBudget, err := b.BudgetsService.GetBudgetWithCategories(budget.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Budget created but failed to load with categories")
	}

	return common.SendSuccessResponse(ctx, "Budget created successfully", createdBudget)
}

func (b *BudgetsController) GetPaginatedBudgets(ctx echo.Context) error {
	userID, ok := ctx.Get("user").(uint)
	if !ok {
		return common.SendInternalServerErrorResponse(ctx, "User authentication required")
	}

	page := 1
	limit := 10
	sort := "created_at DESC"

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

	// if sortParam := ctx.QueryParam("sort"); sortParam != "" {
	// 	// Validate sort parameter to prevent SQL injection
	// 	allowedSorts := map[string]string{
	// 		"created_at": "created_at",
	// 		"updated_at": "updated_at",
	// 		"amount":     "amount",
	// 		"title":      "title",
	// 		"date":       "date",
	// 	}
	// 	if allowedSort, exists := allowedSorts[sortParam]; exists {
	// 		sort = allowedSort + " DESC"
	// 	}
	// }

	budgets, total, err := b.BudgetsService.GetPaginatedBudgetsByUserID(userID, limit, offset, sort)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to fetch budgets")
	}

	response := map[string]interface{}{
		"budgets": budgets,
		"total":   total,
		"page":    page,
		"limit":   limit,
		"sort":    sort,
	}

	return common.SendSuccessResponse(ctx, "Budgets retrieved successfully", response)
}

func (b *BudgetsController) UpdateBudget(ctx echo.Context) error {
	userID, ok := ctx.Get("user").(uint)
	if !ok {
		return common.SendInternalServerErrorResponse(ctx, "User authentication required")
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid budget Id")
	}

	existingBudget, err := b.BudgetsService.GetBudgetByID(uint(id))
	if err != nil {
		return common.SendNotFoundResponse(ctx, "Buget not found")
	}

	if existingBudget.UserID != userID {
		return common.SendUnauthorizedResponse(ctx, "You are not allowed to update this budget")
	}

	request := new(dtos.UpdateBudgetRequestDTO)
	if err := ctx.Bind(request); err != nil {
		return common.SendBadRequestResponse(ctx, err.Error())
	}

	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		return common.SendFailedValidationResponse(ctx, validationErrors)
	}

	if request.Title != "" {
		existingBudget.Title = request.Title
		slug := strings.ToLower(request.Title)
		slug = strings.Replace(slug, " ", "_", -1)
		existingBudget.Slug = slug
	}

	if request.Amount > 0 {
		existingBudget.Amount = request.Amount
	}

	if request.Description != nil {
		existingBudget.Description = request.Description
	}

	if request.Date != "" {
		var parsedDate time.Time
		var parseErr error
		if len(request.Date) == 10 {
			parsedDate, parseErr = time.Parse("2006-01-02", request.Date)
		} else {
			parsedDate, parseErr = time.Parse("2006-01-02 15:04:05", request.Date)
		}
		if parseErr != nil {
			return common.SendBadRequestResponse(ctx, "Invalid date format. Use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS")
		}
		existingBudget.Date = parsedDate
		existingBudget.Month = uint(parsedDate.Month())
		existingBudget.Year = uint16(parsedDate.Year())
	}

	// CheckBudgetsExitExcludingID checks for an existing budget matching the
	// provided user/month/year/slug combination, while excluding the given budget ID.
	exists, err := b.BudgetsService.CheckBudgetsExitExcludingID(existingBudget.UserID, existingBudget.Month, existingBudget.Year, existingBudget.Slug, existingBudget.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, err.Error())
	}
	if exists {
		return common.SendBadRequestResponse(ctx, "Budget with this userId, month, year and slug already exits")
	}

	if err := b.BudgetsService.UpdateBudget(existingBudget); err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to update budget")
	}

	// If categories field is present (even empty), update association accordingly
	if request.Categories != nil {
		var categories []models.Categories
		if len(request.Categories) > 0 {
			for _, catID := range request.Categories {
				if _, err := b.CategoriesService.GetCategoryByID(uint(catID)); err != nil {
					return common.SendBadRequestResponse(ctx, fmt.Sprintf("Category with ID %d does not exist", catID))
				}
			}
			var fetchErr error
			categories, fetchErr = b.CategoriesService.GetCategoriesByIDs(request.Categories)
			if fetchErr != nil {
				return common.SendInternalServerErrorResponse(ctx, "Failed to fetch categories")
			}
		} else {
			categories = []models.Categories{}
		}
		if err := b.BudgetsService.UpdateBudgetCategories(existingBudget.ID, categories); err != nil {
			return common.SendInternalServerErrorResponse(ctx, "Failed to attach categories to budget")
		}
	}

	updatedBudget, err := b.BudgetsService.GetBudgetWithCategories(existingBudget.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Budget updated but failed to load with categories")
	}

	return common.SendSuccessResponse(ctx, "Budget has been updated successfully", updatedBudget)
}
