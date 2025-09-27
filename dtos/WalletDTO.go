package dtos

type CreateWalletRequestDTO struct {
	Name    string  `json:"name" validate:"required"`
	Balance float64 `json:"balance" validate:"required,numeric,min=0"`
}
