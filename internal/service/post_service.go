package service

import (
	"errors"
	"strings"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type PostService struct {
	postRepo    *repository.PostRepository
	userRepo    *repository.UserRepository
	likeRepo    *repository.LikeRepository
	commentRepo *repository.CommentRepository
	followRepo  *repository.FollowRepository
}

func NewPostService(
	postRepo *repository.PostRepository,
	userRepo *repository.UserRepository,
	likeRepo *repository.LikeRepository,
	commentRepo *repository.CommentRepository,
	followRepo *repository.FollowRepository,
) *PostService {
	return &PostService{
		postRepo:    postRepo,
		userRepo:    userRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		followRepo:  followRepo,
	}
}

func (s *PostService) CreatePost(userID string, content string, imageURLs []string, isPublic bool) (*models.Post, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	post := &models.Post{
		UserID:   userID,
		Content:  content,
		IsPublic: isPublic,
	}

	if len(imageURLs) > 0 {
		post.ImageURLs.String = "{" + strings.Join(imageURLs, ",") + "}"
		post.ImageURLs.Valid = true
	}

	if err := s.postRepo.CreatePost(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) UpdatePost(postID, userID string, content *string, imageURLs []string, isPublic *bool) (*models.Post, error) {
	isOwner, err := s.postRepo.IsPostOwner(postID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("unauthorized: you can only update your own posts")
	}

	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	if content != nil {
		if *content == "" {
			return nil, errors.New("content cannot be empty")
		}
		post.Content = *content
	}
	if imageURLs != nil {
		if len(imageURLs) > 0 {
			post.ImageURLs.String = "{" + strings.Join(imageURLs, ",") + "}"
			post.ImageURLs.Valid = true
		} else {
			post.ImageURLs.Valid = false
		}
	}
	if isPublic != nil {
		post.IsPublic = *isPublic
	}

	if err := s.postRepo.UpdatePost(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(postID, userID string) error {
	isOwner, err := s.postRepo.IsPostOwner(postID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("unauthorized: you can only delete your own posts")
	}

	return s.postRepo.DeletePost(postID)
}

func (s *PostService) GetPostByID(postID string) (*models.Post, error) {
	return s.postRepo.GetPostByID(postID)
}

func (s *PostService) GetPosts(limit, offset int) ([]*models.Post, error) {
	return s.postRepo.GetPosts(limit, offset)
}

func (s *PostService) GetUserPosts(userID string, limit, offset int) ([]*models.Post, error) {
	return s.postRepo.GetUserPosts(userID, limit, offset)
}

func (s *PostService) GetFeed(userID string, limit, offset int) ([]*models.Post, error) {
	followingIDs, err := s.followRepo.GetFollowingIDs(userID)
	if err != nil {
		return nil, err
	}

	if len(followingIDs) == 0 {
		return []*models.Post{}, nil
	}

	return s.postRepo.GetFeedPosts(followingIDs, limit, offset)
}
