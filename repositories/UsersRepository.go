package repositories

import (
	"bougette-backend/models"

	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (u *UsersRepository) GetUsers() ([]models.Users, error) {
	var users []models.Users
	err := u.db.Find(&users).Error
	return users, err
}

func (u *UsersRepository) RegisterUser(user *models.Users) error {
	return u.db.Create(user).Error
}

func (u *UsersRepository) FindUserByEmail(email string) (*models.Users, error) {
	var user models.Users

	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
