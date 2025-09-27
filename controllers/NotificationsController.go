package controllers

import (
	"bougette-backend/common"
	"bougette-backend/models"
	"bougette-backend/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

type NotificationsController struct {
	service *services.NotificationsService
}

func NewNotificationsController(service *services.NotificationsService) *NotificationsController {
	return &NotificationsController{service: service}
}

func (c *NotificationsController) CreateNotification(ctx echo.Context) error {
	notification := new(models.Notifications)
	if err := ctx.Bind(notification); err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid notification")
	}

	err := c.service.CreateNotification(notification)
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to create notification")
	}

	return common.SendSuccessResponse(ctx, "Notification created successfully", notification)
}

func (c *NotificationsController) GetNotificationsByUserID(ctx echo.Context) error {
	userID := ctx.Param("user_id")
	id, err := strconv.Atoi(userID)

	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid user ID")
	}

	notifications, err := c.service.GetNotificationsByUserID(uint(id))
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to get notifications")
	}

	return common.SendSuccessResponse(ctx, "Notifications fetched successfully", notifications)
}

func (c *NotificationsController) MarkNotificationAsRead(ctx echo.Context) error {
	notificationID := ctx.Param("notification_id")
	id, err := strconv.Atoi(notificationID)

	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid notification ID")
	}

	err = c.service.MarkNotificationAsRead(uint(id))
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to mark notification as read")
	}

	return common.SendSuccessResponse(ctx, "Notification marked as read successfully", nil)
}

func (c *NotificationsController) MarkAllNotificationsAsRead(ctx echo.Context) error {
	userID := ctx.Param("user_id")
	id, err := strconv.Atoi(userID)

	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid user ID")
	}

	err = c.service.MarkAllNotificationsAsRead(uint(id))
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to mark all notifications as read")
	}

	return common.SendSuccessResponse(ctx, "All notifications marked as read successfully", nil)
}

func (c *NotificationsController) DeleteNotification(ctx echo.Context) error {
	notificationID := ctx.Param("notification_id")
	id, err := strconv.Atoi(notificationID)

	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid notification ID")
	}

	err = c.service.DeleteNotification(uint(id))
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to delete notification")
	}

	return common.SendSuccessResponse(ctx, "Notification deleted successfully", nil)
}

func (c *NotificationsController) ClearAllNotifications(ctx echo.Context) error {
	userID := ctx.Param("user_id")
	id, err := strconv.Atoi(userID)

	if err != nil {
		return common.SendBadRequestResponse(ctx, "Invalid user ID")
	}

	err = c.service.ClearAllNotifications(uint(id))
	if err != nil {
		return common.SendInternalServerErrorResponse(ctx, "Failed to clear all notifications")
	}

	return common.SendSuccessResponse(ctx, "All notifications cleared successfully", nil)
}
