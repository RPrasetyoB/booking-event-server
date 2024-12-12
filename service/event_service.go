package service

import (
	"booking-event-server/dto"
	"booking-event-server/entity"
	errorhandler "booking-event-server/errorHandler"
	"booking-event-server/helper"
	"booking-event-server/repository"
	"time"
)

type EventService interface {
	CreateEvent(req *dto.CreaEventRequest, userId string) (*dto.EventResponse, error)
}

type eventService struct {
	repoEvent repository.EventRepository
	repoDates repository.DatesRepository
}

func NewEventService(eventRepo repository.EventRepository, datesRepo repository.DatesRepository) *eventService {
	return &eventService{
		repoEvent: eventRepo,
		repoDates: datesRepo,
	}
}

func (s *eventService) CreateEvent(req *dto.CreaEventRequest, userId string) (*dto.EventResponse, error) {
	nanoid, _ := helper.GenerateNanoId()
	confirmedDate := (*time.Time)(nil)

	event := entity.Event{
		ID:             nanoid,
		Event_name:     req.Event_name,
		Location:       req.Location,
		User_id:        userId,
		Confirmed_date: confirmedDate,
	}

	createdEvent, err := s.repoEvent.CreateEvent(&event)
	if err != nil {
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	for _, dateStr := range req.Proposed_dates {
		nanoid, _ := helper.GenerateNanoId()
		parsedDate, err := time.Parse("02-01-2006", dateStr)
		if err != nil {
			return nil, &errorhandler.BadRequestError{
				Message: "Invalid date format",
			}
		}
		proposedDate := entity.ProposedDates{
			ID:       nanoid,
			Date:     parsedDate,
			Event_id: createdEvent.ID,
		}
		_, err = s.repoDates.CreateDates(&proposedDate)
		if err != nil {
			return nil, &errorhandler.InternalServerError{
				Message: err.Error(),
			}
		}
	}

	response := &dto.EventResponse{
		ID:             createdEvent.ID,
		Event_name:     createdEvent.Event_name,
		Proposed_dates: req.Proposed_dates,
		Location:       createdEvent.Location,
		User_id:        createdEvent.User_id,
		Confirmed_date: createdEvent.Confirmed_date,
		Created_at:     createdEvent.Created_at,
		Updated_at:     createdEvent.Updated_at,
	}

	return response, nil
}
