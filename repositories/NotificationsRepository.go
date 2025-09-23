package repositories

import (
	"bougette-backend/models"

	"gorm.io/gorm"
)

type NotificationsRepository struct {
	db *gorm.DB
}

func NewNotificationsRepos(db *gorm.DB) *NotificationsRepository {
	return &NotificationsRepository{db: db}
}

func (r *NotificationsRepository) CreateNotification(notification *models.NotificationsModel) error {
	return r.db.Create(notification).Error
}

func (r *NotificationsRepository) GetNotificationsByUserID(userID uint) ([]models.NotificationsModel, error) {
	var notifications []models.NotificationsModel
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationsRepository) MarkNotificationAsRead(id uint) error {
	return r.db.Model(&models.NotificationsModel{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *NotificationsRepository) MarkAllNotificationsAsRead(userID uint) error {
	return r.db.Model(&models.NotificationsModel{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *NotificationsRepository) DeleteNotification(id uint) error {
	return r.db.Delete(&models.NotificationsModel{}, id).Error
}

func (r *NotificationsRepository) ClearAllNotifications(userID uint) error {
	return r.db.Delete(&models.NotificationsModel{}, "user_id = ?", userID).Error
}
