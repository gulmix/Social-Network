package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type StoryService struct {
	storyRepo  *repository.StoryRepository
	followRepo *repository.FollowRepository
}

func NewStoryService(storyRepo *repository.StoryRepository, followRepo *repository.FollowRepository) *StoryService {
	return &StoryService{storyRepo: storyRepo, followRepo: followRepo}
}

func (s *StoryService) CreateStory(userID, mediaURL string, content *string) (*models.Story, error) {
	if mediaURL == "" {
		return nil, errors.New("media URL is required for stories")
	}

	story := &models.Story{
		UserID:    userID,
		MediaURL:  mediaURL,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if content != nil {
		story.Content = sql.NullString{String: *content, Valid: true}
	}

	if err := s.storyRepo.Create(story); err != nil {
		return nil, err
	}
	return story, nil
}

func (s *StoryService) GetStory(storyID string) (*models.Story, error) {
	return s.storyRepo.GetByID(storyID)
}

func (s *StoryService) GetUserStories(userID string) ([]*models.Story, error) {
	return s.storyRepo.GetActiveByUser(userID)
}

func (s *StoryService) GetFeedStories(userID string) ([]*models.Story, error) {
	followingIDs, err := s.followRepo.GetFollowingIDs(userID)
	if err != nil {
		return nil, err
	}
	followingIDs = append(followingIDs, userID)
	return s.storyRepo.GetFeedStories(followingIDs)
}

func (s *StoryService) DeleteStory(storyID, userID string) error {
	isOwner, err := s.storyRepo.IsOwner(storyID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("you can only delete your own stories")
	}
	return s.storyRepo.Delete(storyID)
}

func (s *StoryService) ViewStory(storyID, viewerID string) error {
	return s.storyRepo.AddView(storyID, viewerID)
}

func (s *StoryService) GetStoryViews(storyID, userID string) ([]*models.StoryView, error) {
	isOwner, err := s.storyRepo.IsOwner(storyID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("only the story owner can view story views")
	}
	return s.storyRepo.GetViews(storyID)
}

func (s *StoryService) GetViewCount(storyID string) (int, error) {
	return s.storyRepo.GetViewCount(storyID)
}
