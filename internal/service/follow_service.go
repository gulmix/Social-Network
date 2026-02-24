package service

import (
	"errors"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type FollowService struct {
	followRepo *repository.FollowRepository
	userRepo   *repository.UserRepository
}

func NewFollowService(
	followRepo *repository.FollowRepository,
	userRepo *repository.UserRepository,
) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

func (s *FollowService) FollowUser(followerID, followingID string) (*models.Follow, error) {
	if followerID == followingID {
		return nil, errors.New("cannot follow yourself")
	}

	if _, err := s.userRepo.GetUserByID(followingID); err != nil {
		return nil, errors.New("user not found")
	}

	exists, err := s.followRepo.FollowExists(followerID, followingID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("already following this user")
	}

	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := s.followRepo.CreateFollow(follow); err != nil {
		return nil, err
	}

	return follow, nil
}

func (s *FollowService) UnfollowUser(followerID, followingID string) error {
	exists, err := s.followRepo.FollowExists(followerID, followingID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("not following this user")
	}

	return s.followRepo.DeleteFollow(followerID, followingID)
}

func (s *FollowService) GetFollowers(userID string) ([]*models.User, error) {
	followers, err := s.followRepo.GetFollowers(userID)
	if err != nil {
		return nil, err
	}

	var users []*models.User
	for _, follow := range followers {
		user, err := s.userRepo.GetUserByID(follow.FollowerID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *FollowService) GetFollowing(userID string) ([]*models.User, error) {
	following, err := s.followRepo.GetFollowing(userID)
	if err != nil {
		return nil, err
	}

	var users []*models.User
	for _, follow := range following {
		user, err := s.userRepo.GetUserByID(follow.FollowingID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *FollowService) IsFollowing(followerID, followingID string) (bool, error) {
	return s.followRepo.FollowExists(followerID, followingID)
}

func (s *FollowService) GetFollowersCount(userID string) (int, error) {
	return s.followRepo.GetFollowersCount(userID)
}

func (s *FollowService) GetFollowingCount(userID string) (int, error) {
	return s.followRepo.GetFollowingCount(userID)
}
