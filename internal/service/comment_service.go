package service

import (
	"errors"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
	userRepo    *repository.UserRepository
}

func NewCommentService(
	commentRepo *repository.CommentRepository,
	postRepo *repository.PostRepository,
	userRepo *repository.UserRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) CreateComment(userID, postID, content string, parentID *string) (*models.Comment, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	if _, err := s.postRepo.GetPostByID(postID); err != nil {
		return nil, errors.New("post not found")
	}

	if parentID != nil {
		if _, err := s.commentRepo.GetCommentByID(*parentID); err != nil {
			return nil, errors.New("parent comment not found")
		}
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

	if parentID != nil {
		comment.ParentID.String = *parentID
		comment.ParentID.Valid = true
	}

	if err := s.commentRepo.CreateComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetComment(id string) (*models.Comment, error) {
	return s.commentRepo.GetCommentByID(id)
}

func (s *CommentService) GetCommentByID(id string) (*models.Comment, error) {
	return s.commentRepo.GetCommentByID(id)
}

func (s *CommentService) GetCommentsByPostID(postID string, limit, offset int) ([]*models.Comment, error) {
	return s.commentRepo.GetCommentsByPostID(postID, limit, offset)
}

func (s *CommentService) GetReplies(commentID string) ([]*models.Comment, error) {
	return s.commentRepo.GetRepliesByCommentID(commentID)
}

func (s *CommentService) UpdateComment(commentID, userID, content string) (*models.Comment, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	isOwner, err := s.commentRepo.IsCommentOwner(commentID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("unauthorized: you can only update your own comments")
	}

	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}

	comment.Content = content

	if err := s.commentRepo.UpdateComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) DeleteComment(commentID, userID string) error {
	isOwner, err := s.commentRepo.IsCommentOwner(commentID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("unauthorized: you can only delete your own comments")
	}

	return s.commentRepo.DeleteComment(commentID)
}

func (s *CommentService) GetCommentsCountByPostID(postID string) (int, error) {
	return s.commentRepo.GetCommentsCountByPostID(postID)
}
