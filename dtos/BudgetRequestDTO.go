package dtos

type CreateBudgetRequestDTO struct {
	// dive validate each element individually using the rules you define
	Categories []uint64 `json:"categories" validate:"required,dive,min=1"`
	Amount     float64  `json:"amount" validate:"required,numeric,min=1"`
	// here 2006-01-02 is NOT a literal date - it's Go's special reference format
	Date        string  `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Title       string  `json:"title" validate:"required,min=1,max=250"`
	Description *string `json:"description" validate:"omitempty,min=1,max=500"`
}

type UpdateBudgetRequestDTO struct {
	Categories  []uint64 `json:"categories" validate:"dive,min=1"`
	Amount      float64  `json:"amount" validate:"numeric,min=1"`
	Date        string   `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Title       string   `json:"title" validate:"min=1,max=250"`
	Description *string  `json:"description" validate:"omitempty,min=1,max=500"`
}
