package models

type Notifications struct {
	UserID  uint   `gorm:"not null" json:"user_id"`
	Message string `gorm:"type: text; not null" json:"message"`
	IsRead  bool   `gorm:"type: boolean; not null; default: false" json:"is_read"`
	BaseModel
}
