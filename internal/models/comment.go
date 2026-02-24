package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        string         `json:"id" db:"id"`
	PostID    string         `json:"postId" db:"post_id"`
	UserID    string         `json:"userId" db:"user_id"`
	Content   string         `json:"content" db:"content"`
	ParentID  sql.NullString `json:"parentId" db:"parent_id"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time      `json:"updatedAt" db:"updated_at"`
}

func (Comment) TableName() string {
	return "comments"
}