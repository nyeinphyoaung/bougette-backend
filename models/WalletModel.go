package models

type Wallet struct {
	BaseModel
	UserID  uint    `gorm:"not null; column:user_id;uniqueIndex:unique_user_id_name" json:"user_id"`
	Balance float64 `gorm:"type:double precision;not null;default:0.0" json:"balance"`
	Name    string  `gorm:"type:varchar(255);not null;uniqueIndex:unique_user_id_name" json:"name"`
	Owner   Users   `gorm:"foreignKey:UserID" json:"owner,omitempty"`
}
