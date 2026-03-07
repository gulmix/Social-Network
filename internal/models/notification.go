package models

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID          string         `json:"id" db:"id"`
	UserID      string         `json:"userId" db:"user_id"`
	ActorID     sql.NullString `json:"actorId" db:"actor_id"`
	Type        string         `json:"type" db:"type"`
	ReferenceID sql.NullString `json:"referenceId" db:"reference_id"`
	Content     string         `json:"content" db:"content"`
	IsRead      bool           `json:"isRead" db:"is_read"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
}
