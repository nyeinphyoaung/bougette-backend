package models

import "time"

type Budgets struct {
	Title string `gorm:"index;type:varchar(255);not null" json:"title"`
	// using uniqueIndex:unique_user_id_slug_month_year that is make sure
	// there is no column in our db with the same of user_id, slug, month and year
	Slug        string       `gorm:"index,type:varchar(255);not null;uniqueIndex:unique_user_id_slug_month_year" json:"slug"`
	UserID      uint         `gorm:"not null;column:user_id;uniqueIndex:unique_user_id_slug_month_year" json:"user_id"`
	Description *string      `gorm:"type:text" json:"description"`
	Amount      float64      `gorm:"type:decimal(10,2);not null" json:"amount"`
	Categories  []Categories `gorm:"constraint:OnDelete:CASCADE;many2many:budget_categories" json:"categories"`
	Date        time.Time    `gorm:"type:datetime;not null" json:"date"`
	// by using composite indexes make sure our query faster
	// ensures no duplicate budgets for the same user in the same month/year
	Month uint   `gorm:"type:TINYINT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_month_year" json:"month"`
	Year  uint16 `gorm:"type:INT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_month_year" json:"year"`
	BaseModel
}
