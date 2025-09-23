package services

import (
	"bougette-backend/models"
	"bougette-backend/repositories"
)

type NotificationsService struct {
	repo *repositories.NotificationsRepository
}

func (s *NotificationsService) CreateNotification(notification *models.NotificationsModel) error {
	return s.repo.CreateNotification(notification)
}

func NewNotificationsService(repo *repositories.NotificationsRepository) *NotificationsService {
	return &NotificationsService{repo: repo}
}

func (s *NotificationsService) GetNotificationsByUserID(userID uint) ([]models.NotificationsModel, error) {
	return s.repo.GetNotificationsByUserID(userID)
}

func (s *NotificationsService) MarkNotificationAsRead(id uint) error {
	return s.repo.MarkNotificationAsRead(id)
}

func (s *NotificationsService) MarkAllNotificationsAsRead(userID uint) error {
	return s.repo.MarkAllNotificationsAsRead(userID)
}

func (s *NotificationsService) DeleteNotification(id uint) error {
	return s.repo.DeleteNotification(id)
}

func (s *NotificationsService) ClearAllNotifications(userID uint) error {
	return s.repo.ClearAllNotifications(userID)
}
