package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID          string         `json:"id" db:"id"`
	UserID      string         `json:"userId" db:"user_id"`
	Content     string         `json:"content" db:"content"`
	ImageURLs   sql.NullString `json:"imageUrls" db:"image_urls"`
	IsPublic    bool           `json:"isPublic" db:"is_public"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
}

func (Post) TableName() string {
	return "posts"
}