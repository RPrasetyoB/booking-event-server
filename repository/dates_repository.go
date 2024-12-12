package repository

import (
	"booking-event-server/entity"

	"gorm.io/gorm"
)

type DatesRepository interface {
	CreateDates(dates *entity.ProposedDates) (*entity.ProposedDates, error)
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
