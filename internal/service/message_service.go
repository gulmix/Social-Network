package service

import (
	"database/sql"
	"errors"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type MessageService struct {
	msgRepo  *repository.MessageRepository
	convRepo *repository.ConversationRepository
}

func NewMessageService(msgRepo *repository.MessageRepository, convRepo *repository.ConversationRepository) *MessageService {
	return &MessageService{msgRepo: msgRepo, convRepo: convRepo}
}

func (s *MessageService) SendMessage(senderID, conversationID string, content *string, mediaURL *string) (*models.Message, error) {
	isParticipant, err := s.convRepo.IsParticipant(conversationID, senderID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, errors.New("you are not a participant of this conversation")
	}

	if (content == nil || *content == "") && (mediaURL == nil || *mediaURL == "") {
		return nil, errors.New("message must have content or media")
	}

	msg := &models.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
	}
	if content != nil {
		msg.Content = sql.NullString{String: *content, Valid: true}
	}
	if mediaURL != nil {
		msg.MediaURL = sql.NullString{String: *mediaURL, Valid: true}
	}

	if err := s.msgRepo.Create(msg); err != nil {
		return nil, err
	}

	_ = s.convRepo.UpdateTimestamp(conversationID)

	return msg, nil
}

func (s *MessageService) GetMessages(userID, conversationID string, limit, offset int) ([]*models.Message, error) {
	isParticipant, err := s.convRepo.IsParticipant(conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, errors.New("you are not a participant of this conversation")
	}

	return s.msgRepo.GetByConversation(conversationID, limit, offset)
}

func (s *MessageService) DeleteMessage(messageID, userID string) error {
	isSender, err := s.msgRepo.IsSender(messageID, userID)
	if err != nil {
		return err
	}
	if !isSender {
		return errors.New("you can only delete your own messages")
	}
	return s.msgRepo.Delete(messageID)
}

func (s *MessageService) MarkAsRead(conversationID, userID string) error {
	isParticipant, err := s.convRepo.IsParticipant(conversationID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.New("you are not a participant of this conversation")
	}
	return s.msgRepo.MarkConversationAsRead(conversationID, userID)
}
