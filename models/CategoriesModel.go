package models

type Categories struct {
	Name     string `gorm:"type:varchar(200);unique,not null" json:"name"`
	Slug     string `gorm:"type:varchar(200);unique,not null" json:"slug"`
	IsCustom bool   `gorm:"type:bool;default:false" json:"is_custom"`
	BaseModel
}
