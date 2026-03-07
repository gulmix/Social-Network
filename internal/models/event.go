package models

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          string         `json:"id" db:"id"`
	CreatorID   string         `json:"creatorId" db:"creator_id"`
	GroupID     sql.NullString `json:"groupId" db:"group_id"`
	Title       string         `json:"title" db:"title"`
	Description sql.NullString `json:"description" db:"description"`
	Location    sql.NullString `json:"location" db:"location"`
	StartTime   time.Time      `json:"startTime" db:"start_time"`
	EndTime     sql.NullTime   `json:"endTime" db:"end_time"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
}

type EventAttendee struct {
	ID        string    `json:"id" db:"id"`
	EventID   string    `json:"eventId" db:"event_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
