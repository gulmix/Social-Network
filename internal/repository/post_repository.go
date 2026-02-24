package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) CreatePost(post *models.Post) error {
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	query := `
		INSERT INTO posts (id, user_id, content, image_urls, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		post.ID, post.UserID, post.Content,
		post.ImageURLs, post.IsPublic,
		post.CreatedAt, post.UpdatedAt,
	)
	return err
}

func (r *PostRepository) GetPostByID(id string) (*models.Post, error) {
	query := `
		SELECT id, user_id, content, image_urls, is_public, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	post := &models.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.UserID, &post.Content,
		&post.ImageURLs, &post.IsPublic,
		&post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) GetPosts(limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT id, user_id, content, image_urls, is_public, created_at, updated_at
		FROM posts
		WHERE is_public = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(
			&post.ID, &post.UserID, &post.Content,
			&post.ImageURLs, &post.IsPublic,
			&post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetUserPosts(userID string, limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT id, user_id, content, image_urls, is_public, created_at, updated_at
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(
			&post.ID, &post.UserID, &post.Content,
			&post.ImageURLs, &post.IsPublic,
			&post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetFeedPosts(userIDs []string, limit, offset int) ([]*models.Post, error) {
	if len(userIDs) == 0 {
		return []*models.Post{}, nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(userIDs))
	for i := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	placeholderStr := strings.Join(placeholders, ",")

	query := fmt.Sprintf(`
		SELECT DISTINCT p.id, p.user_id, p.content, p.image_urls, p.is_public, p.created_at, p.updated_at
		FROM posts p
		WHERE p.is_public = true AND p.user_id IN (%s)
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, placeholderStr, len(userIDs)+1, len(userIDs)+2)

	// Build arguments array
	args := make([]interface{}, len(userIDs)+2)
	for i, userID := range userIDs {
		args[i] = userID
	}
	args[len(userIDs)] = limit
	args[len(userIDs)+1] = offset

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(
			&post.ID, &post.UserID, &post.Content,
			&post.ImageURLs, &post.IsPublic,
			&post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) UpdatePost(post *models.Post) error {
	post.UpdatedAt = time.Now()

	query := `
		UPDATE posts
		SET content = $1, image_urls = $2, is_public = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query,
		post.Content, post.ImageURLs, post.IsPublic,
		post.UpdatedAt, post.ID,
	)
	return err
}

func (r *PostRepository) DeletePost(id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PostRepository) IsPostOwner(postID, userID string) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM posts
		WHERE id = $1 AND user_id = $2
	)`

	var isOwner bool
	err := r.db.QueryRow(query, postID, userID).Scan(&isOwner)
	return isOwner, err
}
