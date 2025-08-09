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

func (u *UsersRepository) GetUserByID(id uint) (*models.Users, error) {
	var user models.Users

	err := u.db.First(&user, id).Error
	return &user, err
}

func (u *UsersRepository) RegisterUser(user *models.Users) error {
	return u.db.Create(user).Error
}

func (u *UsersRepository) UpdateUser(user *models.Users) error {
	if err := u.db.Model(user).Updates(user).Error; err != nil {
		return err
	}

	return nil
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

func (u *UsersRepository) DeleteUser(id uint) error {
	return u.db.Unscoped().Delete(&models.Users{}, id).Error
}
