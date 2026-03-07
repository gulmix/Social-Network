package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *models.Event) error {
	event.ID = uuid.New().String()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	query := `INSERT INTO events (id, creator_id, group_id, title, description, location, start_time, end_time, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(query, event.ID, event.CreatorID, event.GroupID, event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.CreatedAt, event.UpdatedAt)
	return err
}

func (r *EventRepository) GetByID(id string) (*models.Event, error) {
	query := `SELECT id, creator_id, group_id, title, description, location, start_time, end_time, created_at, updated_at FROM events WHERE id = $1`
	event := &models.Event{}
	err := r.db.QueryRow(query, id).Scan(&event.ID, &event.CreatorID, &event.GroupID, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *EventRepository) Update(event *models.Event) error {
	event.UpdatedAt = time.Now()
	query := `UPDATE events SET title = $1, description = $2, location = $3, start_time = $4, end_time = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.Exec(query, event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.UpdatedAt, event.ID)
	return err
}

func (r *EventRepository) Delete(id string) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EventRepository) GetUpcoming(limit, offset int) ([]*models.Event, error) {
	query := `SELECT id, creator_id, group_id, title, description, location, start_time, end_time, created_at, updated_at FROM events WHERE start_time > $1 ORDER BY start_time LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, time.Now(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		e := &models.Event{}
		if err := rows.Scan(&e.ID, &e.CreatorID, &e.GroupID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *EventRepository) GetByGroup(groupID string, limit, offset int) ([]*models.Event, error) {
	query := `SELECT id, creator_id, group_id, title, description, location, start_time, end_time, created_at, updated_at FROM events WHERE group_id = $1 ORDER BY start_time LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		e := &models.Event{}
		if err := rows.Scan(&e.ID, &e.CreatorID, &e.GroupID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *EventRepository) IsCreator(eventID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1 AND creator_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, eventID, userID).Scan(&exists)
	return exists, err
}

func (r *EventRepository) AddAttendee(eventID, userID, status string) error {
	query := `INSERT INTO event_attendees (id, event_id, user_id, status, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (event_id, user_id) DO UPDATE SET status = $4`
	_, err := r.db.Exec(query, uuid.New().String(), eventID, userID, status, time.Now())
	return err
}

func (r *EventRepository) RemoveAttendee(eventID, userID string) error {
	query := `DELETE FROM event_attendees WHERE event_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, eventID, userID)
	return err
}

func (r *EventRepository) GetAttendees(eventID string) ([]*models.EventAttendee, error) {
	query := `SELECT id, event_id, user_id, status, created_at FROM event_attendees WHERE event_id = $1 ORDER BY created_at`
	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []*models.EventAttendee
	for rows.Next() {
		a := &models.EventAttendee{}
		if err := rows.Scan(&a.ID, &a.EventID, &a.UserID, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		attendees = append(attendees, a)
	}
	return attendees, rows.Err()
}

func (r *EventRepository) GetAttendeeCount(eventID string) (int, error) {
	query := `SELECT COUNT(*) FROM event_attendees WHERE event_id = $1 AND status = 'going'`
	var count int
	err := r.db.QueryRow(query, eventID).Scan(&count)
	return count, err
}
