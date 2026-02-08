package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID              string         `json:"id" db:"id"`
	Email           string         `json:"email" db:"email"`
	Username        string         `json:"username" db:"username"`
	PasswordHash    sql.NullString `json:"-" db:"password_hash"`
	FirstName       sql.NullString `json:"firstName" db:"first_name"`
	LastName        sql.NullString `json:"lastName" db:"last_name"`
	Bio             sql.NullString `json:"bio" db:"bio"`
	AvatarURL       sql.NullString `json:"avatarUrl" db:"avatar_url"`
	EmailVerified   bool           `json:"emailVerified" db:"email_verified"`
	OAuthProvider   sql.NullString `json:"oauthProvider" db:"oauth_provider"`
	OAuthProviderID sql.NullString `json:"oauthProviderId" db:"oauth_provider_id"`
	CreatedAt       time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time      `json:"updatedAt" db:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
