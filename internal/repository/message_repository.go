package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(msg *models.Message) error {
	msg.ID = uuid.New().String()
	msg.CreatedAt = time.Now()

	query := `INSERT INTO messages (id, conversation_id, sender_id, content, media_url, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, msg.ID, msg.ConversationID, msg.SenderID, msg.Content, msg.MediaURL, msg.CreatedAt)
	return err
}

func (r *MessageRepository) GetByID(id string) (*models.Message, error) {
	query := `SELECT id, conversation_id, sender_id, content, media_url, created_at FROM messages WHERE id = $1`
	msg := &models.Message{}
	err := r.db.QueryRow(query, id).Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.MediaURL, &msg.CreatedAt)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (r *MessageRepository) GetByConversation(conversationID string, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content, media_url, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		msg := &models.Message{}
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.MediaURL, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func (r *MessageRepository) Delete(id string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MessageRepository) IsSender(messageID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM messages WHERE id = $1 AND sender_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, messageID, userID).Scan(&exists)
	return exists, err
}

func (r *MessageRepository) MarkAsRead(messageID, userID string) error {
	query := `INSERT INTO message_reads (id, message_id, user_id, read_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, uuid.New().String(), messageID, userID, time.Now())
	return err
}

func (r *MessageRepository) MarkConversationAsRead(conversationID, userID string) error {
	query := `
		INSERT INTO message_reads (id, message_id, user_id, read_at)
		SELECT $1 || m.id, m.id, $2, $3
		FROM messages m
		WHERE m.conversation_id = $4
		AND m.sender_id != $2
		AND NOT EXISTS (SELECT 1 FROM message_reads mr WHERE mr.message_id = m.id AND mr.user_id = $2)
	`
	_, err := r.db.Exec(query, uuid.New().String()[:8], userID, time.Now(), conversationID)
	return err
}

func (r *MessageRepository) GetUnreadCount(conversationID, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages m
		WHERE m.conversation_id = $1
		AND m.sender_id != $2
		AND NOT EXISTS (SELECT 1 FROM message_reads mr WHERE mr.message_id = m.id AND mr.user_id = $2)
	`
	var count int
	err := r.db.QueryRow(query, conversationID, userID).Scan(&count)
	return count, err
}
