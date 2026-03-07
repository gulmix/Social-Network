package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *models.Notification) error {
	n.ID = uuid.New().String()
	n.CreatedAt = time.Now()

	query := `INSERT INTO notifications (id, user_id, actor_id, type, reference_id, content, is_read, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, n.ID, n.UserID, n.ActorID, n.Type, n.ReferenceID, n.Content, n.IsRead, n.CreatedAt)
	return err
}

func (r *NotificationRepository) GetByUserID(userID string, limit, offset int) ([]*models.Notification, error) {
	query := `SELECT id, user_id, actor_id, type, reference_id, content, is_read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		n := &models.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.ActorID, &n.Type, &n.ReferenceID, &n.Content, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (r *NotificationRepository) MarkAsRead(id string) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(userID string) error {
	query := `UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *NotificationRepository) GetUnreadCount(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *NotificationRepository) Delete(id string) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *NotificationRepository) IsOwner(notificationID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM notifications WHERE id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, notificationID, userID).Scan(&exists)
	return exists, err
}
