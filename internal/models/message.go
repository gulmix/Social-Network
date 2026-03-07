package models

import (
	"database/sql"
	"time"
)

type Message struct {
	ID             string         `json:"id" db:"id"`
	ConversationID string         `json:"conversationId" db:"conversation_id"`
	SenderID       string         `json:"senderId" db:"sender_id"`
	Content        sql.NullString `json:"content" db:"content"`
	MediaURL       sql.NullString `json:"mediaUrl" db:"media_url"`
	CreatedAt      time.Time      `json:"createdAt" db:"created_at"`
}

type MessageRead struct {
	ID        string    `json:"id" db:"id"`
	MessageID string    `json:"messageId" db:"message_id"`
	UserID    string    `json:"userId" db:"user_id"`
	ReadAt    time.Time `json:"readAt" db:"read_at"`
}
