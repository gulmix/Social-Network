package service

import (
	"errors"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type LikeService struct {
	likeRepo *repository.LikeRepository
	postRepo *repository.PostRepository
	userRepo *repository.UserRepository
}

func NewLikeService(
	likeRepo *repository.LikeRepository,
	postRepo *repository.PostRepository,
	userRepo *repository.UserRepository,
) *LikeService {
	return &LikeService{
		likeRepo: likeRepo,
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (s *LikeService) LikePost(userID, postID string) (*models.Like, error) {
	if _, err := s.postRepo.GetPostByID(postID); err != nil {
		return nil, errors.New("post not found")
	}

	exists, err := s.likeRepo.LikeExists(userID, postID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("post already liked")
	}

	like := &models.Like{
		UserID: userID,
		PostID: postID,
	}

	if err := s.likeRepo.CreateLike(like); err != nil {
		return nil, err
	}

	return like, nil
}

func (s *LikeService) UnlikePost(userID, postID string) error {
	exists, err := s.likeRepo.LikeExists(userID, postID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("like not found")
	}

	return s.likeRepo.DeleteLike(userID, postID)
}

func (s *LikeService) GetLikesByPostID(postID string) ([]*models.Like, error) {
	return s.likeRepo.GetLikesByPostID(postID)
}

func (s *LikeService) GetLikesByUserID(userID string) ([]*models.Like, error) {
	return s.likeRepo.GetLikesByUserID(userID)
}

func (s *LikeService) IsLiked(userID, postID string) (bool, error) {
	return s.likeRepo.LikeExists(userID, postID)
}

func (s *LikeService) GetLikesCount(postID string) (int, error) {
	return s.likeRepo.GetLikesCountByPostID(postID)
}
