package controllers

import (
	"bougette-backend/common"
	"bougette-backend/dtos"
	"bougette-backend/models"
	"bougette-backend/services"
	"bougette-backend/validation"

	"fmt"
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
