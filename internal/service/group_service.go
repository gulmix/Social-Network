package service

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type GroupService struct {
	groupRepo *repository.GroupRepository
}

func NewGroupService(groupRepo *repository.GroupRepository) *GroupService {
	return &GroupService{groupRepo: groupRepo}
}

func (s *GroupService) CreateGroup(userID, name string, description *string, isPrivate bool) (*models.Group, error) {
	if name == "" {
		return nil, errors.New("group name is required")
	}

	group := &models.Group{
		Name:      name,
		CreatorID: userID,
		IsPrivate: isPrivate,
	}
	if description != nil {
		group.Description = sql.NullString{String: *description, Valid: true}
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	if err := s.groupRepo.AddMember(group.ID, userID, "admin"); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *GroupService) UpdateGroup(groupID, userID, name string, description *string, isPrivate *bool) (*models.Group, error) {
	isAdmin, err := s.groupRepo.IsAdmin(groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("only admins can update the group")
	}

	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		group.Name = name
	}
	if description != nil {
		group.Description = sql.NullString{String: *description, Valid: true}
	}
	if isPrivate != nil {
		group.IsPrivate = *isPrivate
	}

	if err := s.groupRepo.Update(group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) DeleteGroup(groupID, userID string) error {
	isCreator, err := s.groupRepo.IsCreator(groupID, userID)
	if err != nil {
		return err
	}
	if !isCreator {
		return errors.New("only the creator can delete the group")
	}
	return s.groupRepo.Delete(groupID)
}

func (s *GroupService) GetGroup(groupID string) (*models.Group, error) {
	return s.groupRepo.GetByID(groupID)
}

func (s *GroupService) GetGroups(limit, offset int) ([]*models.Group, error) {
	return s.groupRepo.GetAll(limit, offset)
}

func (s *GroupService) GetUserGroups(userID string) ([]*models.Group, error) {
	return s.groupRepo.GetUserGroups(userID)
}

func (s *GroupService) JoinGroup(groupID, userID string) error {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	if group.IsPrivate {
		return errors.New("this group is private, you need an invitation")
	}
	return s.groupRepo.AddMember(groupID, userID, "member")
}

func (s *GroupService) LeaveGroup(groupID, userID string) error {
	isCreator, err := s.groupRepo.IsCreator(groupID, userID)
	if err != nil {
		return err
	}
	if isCreator {
		return errors.New("creator cannot leave the group, delete it instead")
	}
	return s.groupRepo.RemoveMember(groupID, userID)
}

func (s *GroupService) GetMembers(groupID string) ([]*models.GroupMember, error) {
	return s.groupRepo.GetMembers(groupID)
}

func (s *GroupService) GetMemberCount(groupID string) (int, error) {
	return s.groupRepo.GetMemberCount(groupID)
}

func (s *GroupService) IsMember(groupID, userID string) (bool, error) {
	return s.groupRepo.IsMember(groupID, userID)
}

func (s *GroupService) CreateGroupPost(userID, groupID, content string, imageURLs []string) (*models.GroupPost, error) {
	isMember, err := s.groupRepo.IsMember(groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("you must be a member to post in this group")
	}

	post := &models.GroupPost{
		GroupID: groupID,
		UserID:  userID,
		Content: content,
	}
	if len(imageURLs) > 0 {
		post.ImageURLs = sql.NullString{String: "{" + strings.Join(imageURLs, ",") + "}", Valid: true}
	}

	if err := s.groupRepo.CreatePost(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *GroupService) GetGroupPosts(groupID string, limit, offset int) ([]*models.GroupPost, error) {
	return s.groupRepo.GetPosts(groupID, limit, offset)
}

func (s *GroupService) DeleteGroupPost(postID, userID string) error {
	isOwner, err := s.groupRepo.IsPostOwner(postID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("you can only delete your own posts")
	}
	return s.groupRepo.DeletePost(postID)
}
