package models

import (
	"time"
)

type Like struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	PostID    string    `json:"postId" db:"post_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func (Like) TableName() string {
	return "likes"
}