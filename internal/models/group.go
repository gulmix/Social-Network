package models

import (
	"database/sql"
	"time"
)

type Group struct {
	ID          string         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description" db:"description"`
	AvatarURL   sql.NullString `json:"avatarUrl" db:"avatar_url"`
	CreatorID   string         `json:"creatorId" db:"creator_id"`
	IsPrivate   bool           `json:"isPrivate" db:"is_private"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
}

type GroupMember struct {
	ID       string    `json:"id" db:"id"`
	GroupID  string    `json:"groupId" db:"group_id"`
	UserID   string    `json:"userId" db:"user_id"`
	Role     string    `json:"role" db:"role"`
	JoinedAt time.Time `json:"joinedAt" db:"joined_at"`
}

type GroupPost struct {
	ID        string         `json:"id" db:"id"`
	GroupID   string         `json:"groupId" db:"group_id"`
	UserID    string         `json:"userId" db:"user_id"`
	Content   string         `json:"content" db:"content"`
	ImageURLs sql.NullString `json:"imageUrls" db:"image_urls"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time      `json:"updatedAt" db:"updated_at"`
}
