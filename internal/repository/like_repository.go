package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

func (r *LikeRepository) CreateLike(like *models.Like) error {
	like.ID = uuid.New().String()
	like.CreatedAt = time.Now()

	query := `
		INSERT INTO likes (id, user_id, post_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, post_id) DO NOTHING
	`

	result, err := r.db.Exec(query,
		like.ID, like.UserID, like.PostID, like.CreatedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("like already exists")
	}

	return nil
}

func (r *LikeRepository) GetLikeByID(id string) (*models.Like, error) {
	like := &models.Like{}
	query := `
		SELECT id, user_id, post_id, created_at
		FROM likes
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&like.ID, &like.UserID, &like.PostID, &like.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("like not found")
	}
	if err != nil {
		return nil, err
	}

	return like, nil
}

func (r *LikeRepository) GetLikesByPostID(postID string) ([]*models.Like, error) {
	query := `
		SELECT id, user_id, post_id, created_at
		FROM likes
		WHERE post_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []*models.Like
	for rows.Next() {
		like := &models.Like{}
		err := rows.Scan(
			&like.ID, &like.UserID, &like.PostID, &like.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}

func (r *LikeRepository) GetLikesByUserID(userID string) ([]*models.Like, error) {
	query := `
		SELECT id, user_id, post_id, created_at
		FROM likes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []*models.Like
	for rows.Next() {
		like := &models.Like{}
		err := rows.Scan(
			&like.ID, &like.UserID, &like.PostID, &like.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	return likes, nil
}

func (r *LikeRepository) DeleteLike(userID, postID string) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`
	result, err := r.db.Exec(query, userID, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

func (r *LikeRepository) LikeExists(userID, postID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)`
	err := r.db.QueryRow(query, userID, postID).Scan(&exists)
	return exists, err
}

func (r *LikeRepository) GetLikesCountByPostID(postID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM likes WHERE post_id = $1`
	err := r.db.QueryRow(query, postID).Scan(&count)
	return count, err
}