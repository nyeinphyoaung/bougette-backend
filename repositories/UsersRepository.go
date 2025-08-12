package repositories

import (
	"bougette-backend/models"
	"errors"

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

func (u *UsersRepository) CreatePasswordReset(passwordReset *models.PasswordReset) error {
	return u.db.Create(passwordReset).Error
}

func (u *UsersRepository) GetPasswordResetByTokenAndEmail(token, email string) (*models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	user, err := u.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found for email")
	}

	err = u.db.Where("token = ? AND user_id = ?", token, user.ID).First(&passwordReset).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("password reset token not found or invalid")
	}
	return &passwordReset, err
}

func (u *UsersRepository) InvalidateOldPasswordResets(userID uint) error {
	return u.db.Model(&models.PasswordReset{}).
		Where("user_id = ? AND used = ?", userID, false).
		Update("used", true).Error
}

func (u *UsersRepository) ResetUserPasswordAndMarkTokenUsed(
	user *models.Users,
	newHashedPassword string,
	passwordReset *models.PasswordReset,
) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		user.Password = newHashedPassword
		if err := tx.Model(user).Select("password").Updates(user).Error; err != nil {
			return err
		}

		passwordReset.Used = true
		if err := tx.Save(passwordReset).Error; err != nil {
			return err
		}
		return nil
	})
}
