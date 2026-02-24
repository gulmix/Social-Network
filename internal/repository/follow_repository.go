package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{
		db: db,
	}
}

func (r *FollowRepository) CreateFollow(follow *models.Follow) error {
	follow.ID = uuid.New().String()
	follow.CreatedAt = time.Now()

	query := `
		INSERT INTO follows (id, follower_id, following_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`

	result, err := r.db.Exec(query,
		follow.ID, follow.FollowerID, follow.FollowingID, follow.CreatedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("follow already exists")
	}

	return nil
}

func (r *FollowRepository) GetFollowByID(id string) (*models.Follow, error) {
	query := `
		SELECT id, follower_id, following_id, created_at
		FROM follows
		WHERE id = $1
	`

	follow := &models.Follow{}
	err := r.db.QueryRow(query, id).Scan(
		&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return follow, nil
}

func (r *FollowRepository) GetFollowers(userID string) ([]*models.Follow, error) {
	query := `
		SELECT id, follower_id, following_id, created_at
		FROM follows
		WHERE following_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*models.Follow
	for rows.Next() {
		follow := &models.Follow{}
		err = rows.Scan(
			&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return follows, nil
}

func (r *FollowRepository) GetFollowing(userID string) ([]*models.Follow, error) {
	query := `
		SELECT id, follower_id, following_id, created_at
		FROM follows
		WHERE follower_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*models.Follow
	for rows.Next() {
		follow := &models.Follow{}
		err = rows.Scan(
			&follow.ID, &follow.FollowerID, &follow.FollowingID, &follow.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return follows, nil
}

func (r *FollowRepository) DeleteFollow(followerID, followingID string) error {
	query := `
		DELETE FROM follows
		WHERE follower_id = $1 AND following_id = $2
	`

	result, err := r.db.Exec(query, followerID, followingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("follow not found")
	}

	return nil
}

func (r *FollowRepository) IsFollowing(followerID, followingID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM follows
			WHERE follower_id = $1 AND following_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(query, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *FollowRepository) FollowExists(followerID, followingID string) (bool, error) {
	return r.IsFollowing(followerID, followingID)
}

func (r *FollowRepository) GetFollowersCount(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM follows WHERE following_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *FollowRepository) GetFollowingCount(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *FollowRepository) GetFollowingIDs(userID string) ([]string, error) {
	query := `
		SELECT following_id
		FROM follows
		WHERE follower_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followingIDs []string
	for rows.Next() {
		var followingID string
		err = rows.Scan(&followingID)
		if err != nil {
			return nil, err
		}
		followingIDs = append(followingIDs, followingID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return followingIDs, nil
}
