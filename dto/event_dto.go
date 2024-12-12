package dto

import "time"

type CreaEventRequest struct {
	Event_name     string   `validate:"required" json:"event_name"`
	Proposed_dates []string `validate:"required,dive,datetime=02-01-2006" json:"proposed_dates"`
	Location       string   `validate:"required" json:"location"`
}

type EventResponse struct {
	ID             string     `json:"id"`
	Event_name     string     `json:"event_name"`
	Proposed_dates []string   `validate:"required,dive,datetime=02-01-2006" json:"proposed_dates"`
	Location       string     `json:"location"`
	User_id        string     `json:"user_id"`
	Confirmed_date *time.Time `json:"confirmed_date"`
	Created_at     time.Time  `json:"created_at"`
	Updated_at     time.Time  `json:"updated_at"`
}
