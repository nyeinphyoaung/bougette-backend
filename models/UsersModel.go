package models

type Users struct {
	FirstName *string `gorm:"type:varchar(200)" json:"first_name"`
	LastName  *string `gorm:"type:varchar(200)" json:"last_name"`
	Email     string  `gorm:"type:varchar(200); not null; unique" json:"email"`
	Gender    *string `gorm:"type:varchar(50)" json:"gender"`
	Password  string  `gorm:"type:varchar(200); not null" json:"-"`
	// When you use json:"-" GORM’s .Updates(struct) ignores zero values for non-pointer fields.
	BaseModel
}
