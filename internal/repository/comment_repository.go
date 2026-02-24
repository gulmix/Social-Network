package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) CreateComment(comment *models.Comment) error {
	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	query := `
		INSERT INTO comments (id, post_id, user_id, content, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		comment.ID, comment.PostID, comment.UserID,
		comment.Content, comment.ParentID,
		comment.CreatedAt, comment.UpdatedAt,
	)

	return err
}

func (r *CommentRepository) GetCommentByID(id string) (*models.Comment, error) {
	comment := &models.Comment{}
	query := `
		SELECT id, post_id, user_id, content, parent_id, created_at, updated_at
		FROM comments
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&comment.ID, &comment.PostID, &comment.UserID,
		&comment.Content, &comment.ParentID,
		&comment.CreatedAt, &comment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("comment not found")
	}
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepository) GetCommentsByPostID(postID string, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, user_id, content, parent_id, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID,
			&comment.Content, &comment.ParentID,
			&comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetRepliesByCommentID(commentID string) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, user_id, content, parent_id, created_at, updated_at
		FROM comments
		WHERE parent_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID,
			&comment.Content, &comment.ParentID,
			&comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) UpdateComment(comment *models.Comment) error {
	comment.UpdatedAt = time.Now()
	query := `
		UPDATE comments
		SET content = $2, updated_at = $3
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		comment.ID, comment.Content, comment.UpdatedAt,
	)

	return err
}

func (r *CommentRepository) DeleteComment(id string) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *CommentRepository) CommentExists(id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	return exists, err
}

func (r *CommentRepository) IsCommentOwner(commentID, userID string) (bool, error) {
	var isOwner bool
	query := `SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1 AND user_id = $2)`
	err := r.db.QueryRow(query, commentID, userID).Scan(&isOwner)
	return isOwner, err
}

func (r *CommentRepository) GetCommentsCountByPostID(postID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`
	err := r.db.QueryRow(query, postID).Scan(&count)
	return count, err
}