package models

import (
	"database/sql"
	"time"
)

type Conversation struct {
	ID        string         `json:"id" db:"id"`
	IsGroup   bool           `json:"isGroup" db:"is_group"`
	Name      sql.NullString `json:"name" db:"name"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time      `json:"updatedAt" db:"updated_at"`
}

type ConversationParticipant struct {
	ID             string    `json:"id" db:"id"`
	ConversationID string    `json:"conversationId" db:"conversation_id"`
	UserID         string    `json:"userId" db:"user_id"`
	JoinedAt       time.Time `json:"joinedAt" db:"joined_at"`
}
