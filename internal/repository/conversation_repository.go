package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(conv *models.Conversation) error {
	conv.ID = uuid.New().String()
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()

	query := `INSERT INTO conversations (id, is_group, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, conv.ID, conv.IsGroup, conv.Name, conv.CreatedAt, conv.UpdatedAt)
	return err
}

func (r *ConversationRepository) GetByID(id string) (*models.Conversation, error) {
	query := `SELECT id, is_group, name, created_at, updated_at FROM conversations WHERE id = $1`
	conv := &models.Conversation{}
	err := r.db.QueryRow(query, id).Scan(&conv.ID, &conv.IsGroup, &conv.Name, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

func (r *ConversationRepository) GetUserConversations(userID string, limit, offset int) ([]*models.Conversation, error) {
	query := `
		SELECT c.id, c.is_group, c.name, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1
		ORDER BY c.updated_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []*models.Conversation
	for rows.Next() {
		conv := &models.Conversation{}
		if err := rows.Scan(&conv.ID, &conv.IsGroup, &conv.Name, &conv.CreatedAt, &conv.UpdatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, conv)
	}
	return convs, rows.Err()
}

func (r *ConversationRepository) AddParticipant(conversationID, userID string) error {
	query := `INSERT INTO conversation_participants (id, conversation_id, user_id, joined_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, uuid.New().String(), conversationID, userID, time.Now())
	return err
}

func (r *ConversationRepository) RemoveParticipant(conversationID, userID string) error {
	query := `DELETE FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, conversationID, userID)
	return err
}

func (r *ConversationRepository) GetParticipants(conversationID string) ([]*models.ConversationParticipant, error) {
	query := `SELECT id, conversation_id, user_id, joined_at FROM conversation_participants WHERE conversation_id = $1`
	rows, err := r.db.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*models.ConversationParticipant
	for rows.Next() {
		p := &models.ConversationParticipant{}
		if err := rows.Scan(&p.ID, &p.ConversationID, &p.UserID, &p.JoinedAt); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}

func (r *ConversationRepository) IsParticipant(conversationID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, conversationID, userID).Scan(&exists)
	return exists, err
}

func (r *ConversationRepository) FindDirectConversation(userID1, userID2 string) (*models.Conversation, error) {
	query := `
		SELECT c.id, c.is_group, c.name, c.created_at, c.updated_at
		FROM conversations c
		WHERE c.is_group = false
		AND (SELECT COUNT(*) FROM conversation_participants WHERE conversation_id = c.id) = 2
		AND EXISTS (SELECT 1 FROM conversation_participants WHERE conversation_id = c.id AND user_id = $1)
		AND EXISTS (SELECT 1 FROM conversation_participants WHERE conversation_id = c.id AND user_id = $2)
		LIMIT 1
	`
	conv := &models.Conversation{}
	err := r.db.QueryRow(query, userID1, userID2).Scan(&conv.ID, &conv.IsGroup, &conv.Name, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

func (r *ConversationRepository) UpdateTimestamp(conversationID string) error {
	query := `UPDATE conversations SET updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), conversationID)
	return err
}

func (r *ConversationRepository) GetParticipantIDs(conversationID string) ([]string, error) {
	query := `SELECT user_id FROM conversation_participants WHERE conversation_id = $1`
	rows, err := r.db.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
