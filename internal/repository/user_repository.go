package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, username, password_hash, first_name, last_name, 
		                  bio, avatar_url, email_verified, oauth_provider, oauth_provider_id, 
		                  created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Exec(query,
		user.ID, user.Email, user.Username,
		user.PasswordHash, user.FirstName, user.LastName,
		user.Bio, user.AvatarURL, user.EmailVerified,
		user.OAuthProvider, user.OAuthProviderID,
		user.CreatedAt, user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, 
		       bio, avatar_url, email_verified, oauth_provider, oauth_provider_id,
		       created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL,
		&user.EmailVerified, &user.OAuthProvider, &user.OAuthProviderID,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, 
		       bio, avatar_url, email_verified, oauth_provider, oauth_provider_id,
		       created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL,
		&user.EmailVerified, &user.OAuthProvider, &user.OAuthProviderID,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByOAuthProvider(provider, providerID string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, 
		       bio, avatar_url, email_verified, oauth_provider, oauth_provider_id,
		       created_at, updated_at
		FROM users
		WHERE oauth_provider = $1 AND oauth_provider_id = $2
	`

	err := r.db.QueryRow(query, provider, providerID).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL,
		&user.EmailVerified, &user.OAuthProvider, &user.OAuthProviderID,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, 
		       bio, avatar_url, email_verified, oauth_provider, oauth_provider_id,
		       created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL,
		&user.EmailVerified, &user.OAuthProvider, &user.OAuthProviderID,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, 
		       bio, avatar_url, email_verified, oauth_provider, oauth_provider_id,
		       created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.PasswordHash,
			&user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL,
			&user.EmailVerified, &user.OAuthProvider, &user.OAuthProviderID,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	query := `
		UPDATE users
		SET email = $2, username = $3, password_hash = $4, first_name = $5,
		    last_name = $6, bio = $7, avatar_url = $8, email_verified = $9,
		    oauth_provider = $10, oauth_provider_id = $11, updated_at = $12
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		user.ID, user.Email, user.Username, user.PasswordHash,
		user.FirstName, user.LastName, user.Bio, user.AvatarURL,
		user.EmailVerified, user.OAuthProvider, user.OAuthProviderID,
		user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := r.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}
