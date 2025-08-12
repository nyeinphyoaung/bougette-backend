package models

import "time"

type PasswordReset struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	// Type      string    `json:"type"`
	Token     string    `json:"token"`
	Used      bool      `json:"used"`
	ExpiredAt time.Time `json:"expired_at"`
	BaseModel
}
