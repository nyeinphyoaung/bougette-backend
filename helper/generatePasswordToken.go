package helper

import "github.com/google/uuid"

func GeneratePasswordToken() string {
	return uuid.New().String()
}
