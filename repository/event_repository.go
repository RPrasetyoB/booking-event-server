package repository

import (
	"booking-event-server/entity"
	"time"

	"gorm.io/gorm"
)

type EventRepository interface {
	CreateEvent(event *entity.Event) (*entity.Event, error)
	GetEventByUserID(userID string) ([]*entity.Event, error)
	GetEventByID(eventID string) (*entity.Event, error)
	GetAllEventVendor(name string) ([]*entity.Event, error)
	GetEventByStatus(status string) ([]*entity.Event, error)
	PutEventByID(eventID string, event entity.Event) (*entity.Event, error)
	DeleteEventByID(eventID string) error
	PatchConfirmEventByID(eventID string, confirmedDate time.Time) (*entity.Event, error)
	PatchRejectEventByID(eventID string, remark string) (*entity.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *eventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r eventRepository) CreateEvent(event *entity.Event) (*entity.Event, error) {
	err := r.db.Create(&event).Error
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r eventRepository) GetEventByUserID(userID string) ([]*entity.Event, error) {
	var events []*entity.Event
	err := r.db.Where("user_id = ?", userID).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r eventRepository) GetEventByID(eventID string) (*entity.Event, error) {
	var event entity.Event
	err := r.db.Where("id = ?", eventID).First(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r eventRepository) GetAllEventVendor(name string) ([]*entity.Event, error) {
	var events []*entity.Event
	err := r.db.Where("vendor_name = ?", name).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r eventRepository) GetEventByStatus(status string) ([]*entity.Event, error) {
	var events []*entity.Event
	err := r.db.Where("status = ?", status).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r eventRepository) PutEventByID(eventID string, event entity.Event) (*entity.Event, error) {
	err := r.db.Where("id = ?", eventID).Updates(event).Error
	if err != nil {
		return nil, err
	}

	var updatedEvent entity.Event
	err = r.db.Where("id = ?", eventID).First(&updatedEvent).Error
	if err != nil {
		return nil, err
	}

	return &updatedEvent, nil
}

func (r eventRepository) DeleteEventByID(eventID string) error {
	result := r.db.Where("id = ?", eventID).Delete(&entity.Event{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r eventRepository) PatchConfirmEventByID(eventID string, confirmedDate time.Time) (*entity.Event, error) {
	err := r.db.Model(&entity.Event{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"confirmed_date": confirmedDate,
			"status":         "accepted",
		}).Error
	if err != nil {
		return nil, err
	}

	var updatedEvent entity.Event
	err = r.db.Where("id = ?", eventID).First(&updatedEvent).Error
	if err != nil {
		return nil, err
	}

	return &updatedEvent, nil
}

func (r eventRepository) PatchRejectEventByID(eventID string, remark string) (*entity.Event, error) {
	err := r.db.Model(&entity.Event{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"remark": remark,
			"status": "rejected",
		}).Error
	if err != nil {
		return nil, err
	}

	var updatedEvent entity.Event
	err = r.db.Where("id = ?", eventID).First(&updatedEvent).Error
	if err != nil {
		return nil, err
	}

	return &updatedEvent, nil
}
