package dtos

type CreateCategoryRequestDTO struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryRequestDTO struct {
	Name     *string `json:"name"`
	Slug     *string `json:"slug"`
	IsCustom *bool   `json:"is_custom"`
}
