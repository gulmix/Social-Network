package service

import (
	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type UserService struct {
	userRepo   *repository.UserRepository
	followRepo *repository.FollowRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
	followRepo *repository.FollowRepository,
) *UserService {
	return &UserService{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetUsers() ([]*models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *UserService) GetFollowers(userID string) ([]*models.Follow, error) {
	return s.followRepo.GetFollowers(userID)
}

func (s *UserService) GetFollowing(userID string) ([]*models.Follow, error) {
	return s.followRepo.GetFollowing(userID)
}

func (s *UserService) GetFollowersCount(userID string) (int, error) {
	return s.followRepo.GetFollowersCount(userID)
}

func (s *UserService) GetFollowingCount(userID string) (int, error) {
	return s.followRepo.GetFollowingCount(userID)
}
