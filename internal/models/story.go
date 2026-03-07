package models

import (
	"database/sql"
	"time"
)

type Story struct {
	ID        string         `json:"id" db:"id"`
	UserID    string         `json:"userId" db:"user_id"`
	MediaURL  string         `json:"mediaUrl" db:"media_url"`
	Content   sql.NullString `json:"content" db:"content"`
	ExpiresAt time.Time      `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
}

type StoryView struct {
	ID       string    `json:"id" db:"id"`
	StoryID  string    `json:"storyId" db:"story_id"`
	ViewerID string    `json:"viewerId" db:"viewer_id"`
	ViewedAt time.Time `json:"viewedAt" db:"viewed_at"`
}
