package main

import (
	"bougette-backend/common"
	"bougette-backend/controllers"
	"bougette-backend/dtos"
	"bougette-backend/repositories"
	"bougette-backend/services"
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

var categories = []string{
	"Electronics",
	"Fashion",
	"Garden",
	"Books",
	"Sports",
	"Toys",
}

func SeedCategories(controller *controllers.CategoriesController) error {
	e := echo.New()
	for _, name := range categories {
		reqDto := &dtos.CreateCategoryRequestDTO{Name: name}
		body, _ := json.Marshal(reqDto)
		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/categories")
		ctx.Set("user", "seeder")

		err := controller.CreateCategory(ctx)
		if err != nil {
			if rec.Code != 400 {
				println("Error seeding category:", name, "-", err.Error())
			}
		}
	}
	return nil
}

func main() {
	db, err := common.GetDB()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	categoriesRepo := repositories.NewCategoriesRepository(db)
	categoriesService := services.NewCategoriesService(categoriesRepo)
	categoriesController := controllers.NewCategoriesController(categoriesService)

	if err := SeedCategories(categoriesController); err != nil {
		panic("failed to seed categories: " + err.Error())
	}
	println("Categories seeded successfully!")
}
