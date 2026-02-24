package models

import (
	"time"
)

type Follow struct {
	ID          string    `json:"id" db:"id"`
	FollowerID  string    `json:"followerId" db:"follower_id"`
	FollowingID string    `json:"followingId" db:"following_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

func (Follow) TableName() string {
	return "follows"
}