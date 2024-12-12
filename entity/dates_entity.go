package entity

import "time"

type ProposedDates struct {
	ID       string    `gorm:"column:id"`
	Date     time.Time `gorm:"type:date"`
	Event_id string
}

func (ProposedDates) TableName() string {
	return "Proposed_dates"
}
