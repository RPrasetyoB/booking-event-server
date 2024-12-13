package repository

import (
	"booking-event-server/entity"

	"gorm.io/gorm"
)

type DatesRepository interface {
	CreateDates(dates *entity.ProposedDates) (*entity.ProposedDates, error)
	GetDatesByEventID(eventID string) ([]*entity.ProposedDates, error)
	DeleteDatesByEventID(eventID string) error
}

type datesRepository struct {
	db *gorm.DB
}

func NewDatesRepository(db *gorm.DB) *datesRepository {
	return &datesRepository{
		db: db,
	}
}

func (r datesRepository) CreateDates(dates *entity.ProposedDates) (*entity.ProposedDates, error) {
	err := r.db.Create(&dates).Error
	if err != nil {
		return nil, err
	}
	return dates, nil
}

func (r datesRepository) GetDatesByEventID(eventID string) ([]*entity.ProposedDates, error) {
	var dates []*entity.ProposedDates
	err := r.db.Where("event_id = ?", eventID).Find(&dates).Error
	if err != nil {
		return nil, err
	}
	return dates, nil
}

func (r datesRepository) DeleteDatesByEventID(eventID string) error {
	result := r.db.Where("event_id = ?", eventID).Delete(&entity.ProposedDates{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
