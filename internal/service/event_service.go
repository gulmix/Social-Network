package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
)

type EventService struct {
	eventRepo *repository.EventRepository
}

func NewEventService(eventRepo *repository.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

func (s *EventService) CreateEvent(userID, title string, description, location *string, groupID *string, startTime time.Time, endTime *time.Time) (*models.Event, error) {
	if title == "" {
		return nil, errors.New("event title is required")
	}
	if startTime.Before(time.Now()) {
		return nil, errors.New("start time must be in the future")
	}

	event := &models.Event{
		CreatorID: userID,
		Title:     title,
		StartTime: startTime,
	}
	if description != nil {
		event.Description = sql.NullString{String: *description, Valid: true}
	}
	if location != nil {
		event.Location = sql.NullString{String: *location, Valid: true}
	}
	if groupID != nil {
		event.GroupID = sql.NullString{String: *groupID, Valid: true}
	}
	if endTime != nil {
		event.EndTime = sql.NullTime{Time: *endTime, Valid: true}
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	if err := s.eventRepo.AddAttendee(event.ID, userID, "going"); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) UpdateEvent(eventID, userID, title string, description, location *string, startTime *time.Time, endTime *time.Time) (*models.Event, error) {
	isCreator, err := s.eventRepo.IsCreator(eventID, userID)
	if err != nil {
		return nil, err
	}
	if !isCreator {
		return nil, errors.New("only the creator can update the event")
	}

	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		event.Title = title
	}
	if description != nil {
		event.Description = sql.NullString{String: *description, Valid: true}
	}
	if location != nil {
		event.Location = sql.NullString{String: *location, Valid: true}
	}
	if startTime != nil {
		event.StartTime = *startTime
	}
	if endTime != nil {
		event.EndTime = sql.NullTime{Time: *endTime, Valid: true}
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventService) DeleteEvent(eventID, userID string) error {
	isCreator, err := s.eventRepo.IsCreator(eventID, userID)
	if err != nil {
		return err
	}
	if !isCreator {
		return errors.New("only the creator can delete the event")
	}
	return s.eventRepo.Delete(eventID)
}

func (s *EventService) GetEvent(eventID string) (*models.Event, error) {
	return s.eventRepo.GetByID(eventID)
}

func (s *EventService) GetUpcomingEvents(limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.GetUpcoming(limit, offset)
}

func (s *EventService) GetGroupEvents(groupID string, limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.GetByGroup(groupID, limit, offset)
}

func (s *EventService) RespondToEvent(eventID, userID, status string) error {
	validStatuses := map[string]bool{"going": true, "interested": true, "not_going": true}
	if !validStatuses[status] {
		return errors.New("invalid status, must be 'going', 'interested', or 'not_going'")
	}
	return s.eventRepo.AddAttendee(eventID, userID, status)
}

func (s *EventService) CancelAttendance(eventID, userID string) error {
	return s.eventRepo.RemoveAttendee(eventID, userID)
}

func (s *EventService) GetAttendees(eventID string) ([]*models.EventAttendee, error) {
	return s.eventRepo.GetAttendees(eventID)
}

func (s *EventService) GetAttendeeCount(eventID string) (int, error) {
	return s.eventRepo.GetAttendeeCount(eventID)
}
