package service

import (
	"booking-event-server/dto"
	"booking-event-server/entity"
	errorhandler "booking-event-server/errorHandler"
	"booking-event-server/helper"
	"booking-event-server/repository"
	"time"

	"gorm.io/gorm"
)

type EventService interface {
	CreateEvent(req *dto.CreaEventRequest, userId string) (*dto.EventResponse, error)
	GetAllEventsHRByUserID(userID string) ([]*dto.GetEventResponse, error)
	GetEventByID(eventID string) (*dto.GetEventResponse, error)
	UpdateEventHR(req *dto.CreaEventRequest, userID string, eventID string) (*dto.EventResponse, error)
	DeleteEvent(eventID string) error
	GetAllEvents() ([]*dto.GetEventResponse, error)
	AcceptEventVendor(req *dto.ConfirmDateRequest, eventID string) (*dto.EventResponse, error)
	RejectEventVendor(req *dto.RejectRequest, eventID string) (*dto.EventResponse, error)
}

type eventService struct {
	repoEvent repository.EventRepository
	repoDates repository.DatesRepository
	repoUser  repository.AuthRepository
}

func NewEventService(
	eventRepo repository.EventRepository,
	datesRepo repository.DatesRepository,
	userRepo repository.AuthRepository) *eventService {
	return &eventService{
		repoEvent: eventRepo,
		repoDates: datesRepo,
		repoUser:  userRepo,
	}
}

func (s *eventService) CreateEvent(req *dto.CreaEventRequest, userID string) (*dto.EventResponse, error) {
	nanoid, _ := helper.GenerateNanoId()
	confirmedDate := (*time.Time)(nil)

	event := entity.Event{
		ID:             nanoid,
		Event_name:     req.Event_name,
		Location:       req.Location,
		User_id:        userID,
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
		Status:         createdEvent.Status,
		Confirmed_date: createdEvent.Confirmed_date,
		Created_at:     createdEvent.Created_at,
		Updated_at:     createdEvent.Updated_at,
	}

	return response, nil
}

func (s *eventService) GetAllEventsHRByUserID(userID string) ([]*dto.GetEventResponse, error) {
	events, err := s.repoEvent.GetEventByUserID(userID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}
	user, err := s.repoUser.FindUserById(userID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	var response []*dto.GetEventResponse
	for _, event := range events {
		dates, err := s.repoDates.GetDatesByEventID(event.ID)
		if err != nil {
			return nil, &errorhandler.InternalServerError{
				Message: err.Error(),
			}
		}

		var dateStrings []string
		for _, date := range dates {
			dateStrings = append(dateStrings, date.Date.Format("02-01-2006"))
		}

		eventResponse := &dto.GetEventResponse{
			ID:             event.ID,
			Event_name:     event.Event_name,
			Proposed_dates: dateStrings,
			Location:       event.Location,
			Status:         event.Status,
			User_id:        event.User_id,
			User_name:      &user.Name,
			Confirmed_date: event.Confirmed_date,
			Created_at:     event.Created_at,
			Updated_at:     event.Updated_at,
		}
		response = append(response, eventResponse)
	}

	return response, nil
}

func (s *eventService) GetEventByID(eventID string) (*dto.GetEventResponse, error) {
	event, err := s.repoEvent.GetEventByID(eventID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}
	user, err := s.repoUser.FindUserById(event.User_id)
	if err != nil {
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	dates, err := s.repoDates.GetDatesByEventID(event.ID)
	if err != nil {
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	var dateStrings []string
	for _, date := range dates {
		dateStrings = append(dateStrings, date.Date.Format("02-01-2006"))
	}

	eventResponse := &dto.GetEventResponse{
		ID:             event.ID,
		Event_name:     event.Event_name,
		Proposed_dates: dateStrings,
		Location:       event.Location,
		Status:         event.Status,
		User_id:        event.User_id,
		User_name:      &user.Name,
		Confirmed_date: event.Confirmed_date,
		Created_at:     event.Created_at,
		Updated_at:     event.Updated_at,
	}

	return eventResponse, nil
}

func (s *eventService) UpdateEventHR(req *dto.CreaEventRequest, userID string, eventID string) (*dto.EventResponse, error) {
	confirmedDate := (*time.Time)(nil)
	event := entity.Event{
		Event_name:     req.Event_name,
		Location:       req.Location,
		User_id:        userID,
		Confirmed_date: confirmedDate,
	}

	getDates, err := s.repoDates.GetDatesByEventID(eventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errorhandler.NotFoundError{
				Message: "Event not found",
			}
		}
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	for _, date := range getDates {
		err := s.repoDates.DeleteDatesByEventID(date.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, &errorhandler.NotFoundError{
					Message: "Proposed date not found",
				}
			}
			return nil, &errorhandler.InternalServerError{
				Message: err.Error(),
			}
		}
	}

	updatedEvent, err := s.repoEvent.PutEventByID(eventID, event)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errorhandler.NotFoundError{
				Message: "Event not found",
			}
		}
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
			Event_id: updatedEvent.ID,
		}
		_, err = s.repoDates.CreateDates(&proposedDate)
		if err != nil {
			return nil, &errorhandler.InternalServerError{
				Message: err.Error(),
			}
		}
	}

	response := &dto.EventResponse{
		ID:             updatedEvent.ID,
		Event_name:     updatedEvent.Event_name,
		Proposed_dates: req.Proposed_dates,
		Location:       updatedEvent.Location,
		Status:         updatedEvent.Status,
		User_id:        updatedEvent.User_id,
		Confirmed_date: updatedEvent.Confirmed_date,
		Created_at:     updatedEvent.Created_at,
		Updated_at:     updatedEvent.Updated_at,
	}

	return response, nil
}

func (s *eventService) DeleteEvent(eventID string) error {
	err := s.repoDates.DeleteDatesByEventID(eventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errorhandler.NotFoundError{
				Message: "Event not found",
			}
		}
		return &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	err2 := s.repoEvent.DeleteEventByID(eventID)
	if err2 != nil {
		if err2 == gorm.ErrRecordNotFound {
			return &errorhandler.NotFoundError{
				Message: "Event not found",
			}
		}
		return &errorhandler.InternalServerError{
			Message: err2.Error(),
		}
	}

	return nil
}

func (s *eventService) GetAllEvents() ([]*dto.GetEventResponse, error) {
	events, err := s.repoEvent.GetAllEvent()
	if err != nil {
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	var response []*dto.GetEventResponse
	for _, event := range events {
		dates, err := s.repoDates.GetDatesByEventID(event.ID)
		if err != nil {
			return nil, &errorhandler.InternalServerError{
				Message: err.Error(),
			}
		}

		user, err := s.repoUser.FindUserById(event.User_id)
		if err != nil {
			return nil, &errorhandler.InternalServerError{
				Message: err.Error(),
			}
		}

		var dateStrings []string
		for _, date := range dates {
			dateStrings = append(dateStrings, date.Date.Format("02-01-2006"))
		}

		eventResponse := &dto.GetEventResponse{
			ID:             event.ID,
			Event_name:     event.Event_name,
			Proposed_dates: dateStrings,
			Location:       event.Location,
			Status:         event.Status,
			User_id:        event.User_id,
			User_name:      &user.Name,
			Confirmed_date: event.Confirmed_date,
			Created_at:     event.Created_at,
			Updated_at:     event.Updated_at,
		}
		response = append(response, eventResponse)
	}

	return response, nil
}

func (s *eventService) AcceptEventVendor(req *dto.ConfirmDateRequest, eventID string) (*dto.EventResponse, error) {
	parsedDate, err := time.Parse("02-01-2006", req.Confirmed_date)
	if err != nil {
		return nil, &errorhandler.InternalServerError{
			Message: "Invalid date format",
		}
	}

	updatedEvent, err := s.repoEvent.PatchConfirmEventByID(eventID, parsedDate)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errorhandler.NotFoundError{
				Message: "Event not found",
			}
		}
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	dates, err := s.repoDates.GetDatesByEventID(eventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errorhandler.NotFoundError{
				Message: "Event not found",
			}
		}
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	err = s.repoDates.DeleteDatesByEventID(eventID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, &errorhandler.InternalServerError{
			Message: "Failed to delete existing proposed dates",
		}
	}

	var dateStrings []string
	for _, date := range dates {
		nanoid, _ := helper.GenerateNanoId()
		proposedDate := entity.ProposedDates{
			ID:       nanoid,
			Date:     date.Date,
			Event_id: eventID,
		}
		newDate, err := s.repoDates.CreateDates(&proposedDate)
		if err != nil {
			return nil, &errorhandler.InternalServerError{
				Message: err.Error(),
			}
		}
		dateStrings = append(dateStrings, newDate.Date.Format("02-01-2006"))
	}

	response := &dto.EventResponse{
		ID:             updatedEvent.ID,
		Event_name:     updatedEvent.Event_name,
		Proposed_dates: dateStrings,
		Location:       updatedEvent.Location,
		Status:         updatedEvent.Status,
		User_id:        updatedEvent.User_id,
		Confirmed_date: updatedEvent.Confirmed_date,
		Created_at:     updatedEvent.Created_at,
		Updated_at:     updatedEvent.Updated_at,
	}

	return response, nil
}

func (s *eventService) RejectEventVendor(req *dto.RejectRequest, eventID string) (*dto.EventResponse, error) {
	getEvent, err := s.repoEvent.GetEventByID(eventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errorhandler.NotFoundError{
				Message: "Event not found",
			}
		}
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}
	if getEvent.Confirmed_date != nil {
		return nil, &errorhandler.BadRequestError{
			Message: "Event's date accepted already, cannot be reject",
		}
	}

	updatedEvent, _ := s.repoEvent.PatchRejectEventByID(eventID, req.Remark)

	dates, err := s.repoDates.GetDatesByEventID(eventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errorhandler.NotFoundError{
				Message: "Event not found",
			}
		}
		return nil, &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	err = s.repoDates.DeleteDatesByEventID(eventID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, &errorhandler.InternalServerError{
			Message: "Failed to delete existing proposed dates",
		}
	}

	var dateStrings []string
	for _, date := range dates {
		nanoid, _ := helper.GenerateNanoId()
		proposedDate := entity.ProposedDates{
			ID:       nanoid,
			Date:     date.Date,
			Event_id: eventID,
		}
		newDate, err := s.repoDates.CreateDates(&proposedDate)
		if err != nil {
			return nil, &errorhandler.InternalServerError{
				Message: err.Error(),
			}
		}
		dateStrings = append(dateStrings, newDate.Date.Format("02-01-2006"))
	}

	response := &dto.EventResponse{
		ID:             updatedEvent.ID,
		Event_name:     updatedEvent.Event_name,
		Proposed_dates: dateStrings,
		Location:       updatedEvent.Location,
		Status:         updatedEvent.Status,
		Remark:         updatedEvent.Remark,
		User_id:        updatedEvent.User_id,
		Confirmed_date: updatedEvent.Confirmed_date,
		Created_at:     updatedEvent.Created_at,
		Updated_at:     updatedEvent.Updated_at,
	}

	return response, nil
}
