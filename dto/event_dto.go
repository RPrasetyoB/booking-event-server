package dto

import "time"

type CreaEventRequest struct {
	Event_name     string   `validate:"required" json:"event_name"`
	Proposed_dates []string `validate:"required,dive,datetime=02-01-2006" json:"proposed_dates"`
	Vendor_name    string   `validate:"required" json:"vendor_name"`
	Location       string   `validate:"required" json:"location"`
}

type ConfirmDateRequest struct {
	Confirmed_date string `json:"confirmed_date"`
}

type RejectRequest struct {
	Remark string `json:"remark"`
}

type EventResponse struct {
	ID             string     `json:"id"`
	Event_name     string     `json:"event_name"`
	Proposed_dates []string   `validate:"required,dive,datetime=02-01-2006" json:"proposed_dates"`
	Vendor_name    string     `json:"vendor_name"`
	Location       string     `json:"location"`
	Status         string     `json:"status"`
	User_id        string     `json:"user_id"`
	Remark         *string    `json:"remark"`
	Confirmed_date *time.Time `json:"confirmed_date"`
	Created_at     time.Time  `json:"created_at"`
	Updated_at     time.Time  `json:"updated_at"`
}

type GetEventResponse struct {
	ID             string     `json:"id"`
	Event_name     string     `json:"event_name"`
	Proposed_dates []string   `validate:"required,dive,datetime=02-01-2006" json:"proposed_dates"`
	Vendor_name    string     `json:"vendor_name"`
	Location       string     `json:"location"`
	Status         string     `json:"status"`
	User_id        string     `json:"user_id"`
	User_name      *string    `json:"user_name"`
	Remark         *string    `json:"remark"`
	Confirmed_date *time.Time `json:"confirmed_date"`
	Created_at     time.Time  `json:"created_at"`
	Updated_at     time.Time  `json:"updated_at"`
}
