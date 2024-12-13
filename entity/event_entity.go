package entity

import "time"

type Event struct {
	ID             string `gorm:"column:id"`
	Event_name     string
	Location       string
	User_id        string
	Status         string `gorm:"default:'pending'"`
	Remark         *string
	Confirmed_date *time.Time `gorm:"type:date"`
	Created_at     time.Time  `gorm:"type:timestamp;autoCreateTime"`
	Updated_at     time.Time  `gorm:"type:timestamp;autoUpdateTime"`
}

func (Event) TableName() string {
	return "Events"
}
