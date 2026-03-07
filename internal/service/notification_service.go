package service

import (
	"database/sql"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type NotificationService struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}

func (s *NotificationService) CreateNotification(userID, actorID, notifType, referenceID, content string) (*models.Notification, error) {
	n := &models.Notification{
		UserID:  userID,
		Type:    notifType,
		Content: content,
	}
	if actorID != "" {
		n.ActorID = sql.NullString{String: actorID, Valid: true}
	}
	if referenceID != "" {
		n.ReferenceID = sql.NullString{String: referenceID, Valid: true}
	}

	if err := s.notifRepo.Create(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *NotificationService) GetNotifications(userID string, limit, offset int) ([]*models.Notification, error) {
	return s.notifRepo.GetByUserID(userID, limit, offset)
}

func (s *NotificationService) MarkAsRead(notificationID, userID string) error {
	isOwner, err := s.notifRepo.IsOwner(notificationID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return err
	}
	return s.notifRepo.MarkAsRead(notificationID)
}

func (s *NotificationService) MarkAllAsRead(userID string) error {
	return s.notifRepo.MarkAllAsRead(userID)
}

func (s *NotificationService) GetUnreadCount(userID string) (int, error) {
	return s.notifRepo.GetUnreadCount(userID)
}
