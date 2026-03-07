package service

import (
	"database/sql"
	"errors"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type ConversationService struct {
	convRepo *repository.ConversationRepository
	userRepo *repository.UserRepository
}

func NewConversationService(convRepo *repository.ConversationRepository, userRepo *repository.UserRepository) *ConversationService {
	return &ConversationService{convRepo: convRepo, userRepo: userRepo}
}

func (s *ConversationService) CreateDirectConversation(userID, otherUserID string) (*models.Conversation, error) {
	if userID == otherUserID {
		return nil, errors.New("cannot create conversation with yourself")
	}

	existing, err := s.convRepo.FindDirectConversation(userID, otherUserID)
	if err == nil && existing != nil {
		return existing, nil
	}

	conv := &models.Conversation{IsGroup: false}
	if err := s.convRepo.Create(conv); err != nil {
		return nil, err
	}

	if err := s.convRepo.AddParticipant(conv.ID, userID); err != nil {
		return nil, err
	}
	if err := s.convRepo.AddParticipant(conv.ID, otherUserID); err != nil {
		return nil, err
	}

	return conv, nil
}

func (s *ConversationService) CreateGroupConversation(userID string, name string, participantIDs []string) (*models.Conversation, error) {
	if name == "" {
		return nil, errors.New("group conversation name is required")
	}

	conv := &models.Conversation{
		IsGroup: true,
		Name:    sql.NullString{String: name, Valid: true},
	}
	if err := s.convRepo.Create(conv); err != nil {
		return nil, err
	}

	if err := s.convRepo.AddParticipant(conv.ID, userID); err != nil {
		return nil, err
	}
	for _, pid := range participantIDs {
		if pid != userID {
			if err := s.convRepo.AddParticipant(conv.ID, pid); err != nil {
				return nil, err
			}
		}
	}

	return conv, nil
}

func (s *ConversationService) GetUserConversations(userID string, limit, offset int) ([]*models.Conversation, error) {
	return s.convRepo.GetUserConversations(userID, limit, offset)
}

func (s *ConversationService) GetConversation(conversationID, userID string) (*models.Conversation, error) {
	isParticipant, err := s.convRepo.IsParticipant(conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, errors.New("you are not a participant of this conversation")
	}
	return s.convRepo.GetByID(conversationID)
}

func (s *ConversationService) GetParticipantIDs(conversationID string) ([]string, error) {
	return s.convRepo.GetParticipantIDs(conversationID)
}
