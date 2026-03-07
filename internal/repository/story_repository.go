package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type StoryRepository struct {
	db *sql.DB
}

func NewStoryRepository(db *sql.DB) *StoryRepository {
	return &StoryRepository{db: db}
}

func (r *StoryRepository) Create(story *models.Story) error {
	story.ID = uuid.New().String()
	story.CreatedAt = time.Now()

	query := `INSERT INTO stories (id, user_id, media_url, content, expires_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, story.ID, story.UserID, story.MediaURL, story.Content, story.ExpiresAt, story.CreatedAt)
	return err
}

func (r *StoryRepository) GetByID(id string) (*models.Story, error) {
	query := `SELECT id, user_id, media_url, content, expires_at, created_at FROM stories WHERE id = $1`
	story := &models.Story{}
	err := r.db.QueryRow(query, id).Scan(&story.ID, &story.UserID, &story.MediaURL, &story.Content, &story.ExpiresAt, &story.CreatedAt)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func (r *StoryRepository) GetActiveByUser(userID string) ([]*models.Story, error) {
	query := `SELECT id, user_id, media_url, content, expires_at, created_at FROM stories WHERE user_id = $1 AND expires_at > $2 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []*models.Story
	for rows.Next() {
		s := &models.Story{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.MediaURL, &s.Content, &s.ExpiresAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		stories = append(stories, s)
	}
	return stories, rows.Err()
}

func (r *StoryRepository) GetFeedStories(userIDs []string) ([]*models.Story, error) {
	if len(userIDs) == 0 {
		return []*models.Story{}, nil
	}

	query := `SELECT id, user_id, media_url, content, expires_at, created_at FROM stories WHERE user_id = ANY($1) AND expires_at > $2 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, pqStringArray(userIDs), time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []*models.Story
	for rows.Next() {
		s := &models.Story{}
		if err := rows.Scan(&s.ID, &s.UserID, &s.MediaURL, &s.Content, &s.ExpiresAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		stories = append(stories, s)
	}
	return stories, rows.Err()
}

func (r *StoryRepository) Delete(id string) error {
	query := `DELETE FROM stories WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *StoryRepository) IsOwner(storyID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM stories WHERE id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, storyID, userID).Scan(&exists)
	return exists, err
}

func (r *StoryRepository) AddView(storyID, viewerID string) error {
	query := `INSERT INTO story_views (id, story_id, viewer_id, viewed_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, uuid.New().String(), storyID, viewerID, time.Now())
	return err
}

func (r *StoryRepository) GetViews(storyID string) ([]*models.StoryView, error) {
	query := `SELECT id, story_id, viewer_id, viewed_at FROM story_views WHERE story_id = $1 ORDER BY viewed_at DESC`
	rows, err := r.db.Query(query, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []*models.StoryView
	for rows.Next() {
		v := &models.StoryView{}
		if err := rows.Scan(&v.ID, &v.StoryID, &v.ViewerID, &v.ViewedAt); err != nil {
			return nil, err
		}
		views = append(views, v)
	}
	return views, rows.Err()
}

func (r *StoryRepository) GetViewCount(storyID string) (int, error) {
	query := `SELECT COUNT(*) FROM story_views WHERE story_id = $1`
	var count int
	err := r.db.QueryRow(query, storyID).Scan(&count)
	return count, err
}

func (r *StoryRepository) DeleteExpired() error {
	query := `DELETE FROM stories WHERE expires_at < $1`
	_, err := r.db.Exec(query, time.Now())
	return err
}

func pqStringArray(arr []string) string {
	result := "{"
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += "\"" + s + "\""
	}
	result += "}"
	return result
}
