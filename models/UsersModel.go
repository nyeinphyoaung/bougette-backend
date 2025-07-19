package models

import "gorm.io/gorm"

type Users struct {
	FirstName *string `gorm:"type:varchar(200)"`
	LastName  *string `gorm:"type:varchar(200)"`
	Email     string  `gorm:"type:varchar(200); not null; unique"`
	Gender    *string `gorm:"type:varchar(50)"`
	Password  string  `gorm:"type:varchar(200); not null"`
	gorm.Model
}
