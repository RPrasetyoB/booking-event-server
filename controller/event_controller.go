package controller

import (
	"booking-event-server/dto"
	errorhandler "booking-event-server/errorHandler"
	"booking-event-server/helper"
	"booking-event-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type eventController struct {
	serviceEvent service.EventService
	serviceAuth  service.AuthService
}

func NewEventController(eventService service.EventService, authService service.AuthService) *eventController {
	return &eventController{
		serviceEvent: eventService,
		serviceAuth:  authService,
	}
}

func (e *eventController) CreateEvent(c *gin.Context) {
	var event dto.CreaEventRequest

	if err := c.ShouldBindJSON(&event); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: "Payload type invalid",
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(event)
	if err != nil {
		errorMsg := helper.GetErrorMessage(err)
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: errorMsg,
		})
		return
	}

	userData, exists := c.Get("user")
	if !exists {
		errorhandler.HandleError(c, &errorhandler.UnauthorizedError{
			Message: "User not found",
		})
		return
	}

	userID, ok := userData.(map[string]interface{})["user_id"].(string)
	if !ok || userID == "" {
		errorhandler.HandleError(c, &errorhandler.UnauthorizedError{
			Message: "Invalid user ID",
		})
		return
	}

	newEvent, err := e.serviceEvent.CreateEvent(&event, userID)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Event created Successfully",
		Data:       newEvent,
	})

	c.JSON(http.StatusOK, res)
}

func (e *eventController) GetEventsHRbyUserID(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		errorhandler.HandleError(c, &errorhandler.UnauthorizedError{
			Message: "User not found",
		})
		return
	}

	userID, _ := userData.(map[string]interface{})["user_id"].(string)
	events, err := e.serviceEvent.GetAllEventsHRByUserID(userID)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	var statusCode int
	var msg string

	if len(events) == 0 {
		statusCode = http.StatusNotFound
		msg = "No Event"
	} else {
		statusCode = http.StatusOK
		msg = "Events retrieved successfully"
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: statusCode,
		Message:    msg,
		Data:       events,
	})

	c.JSON(statusCode, res)
}

func (e *eventController) GetEventbyID(c *gin.Context) {
	eventID := c.Param("id")
	event, err := e.serviceEvent.GetEventByID(eventID)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Event retrieved successfully",
		Data:       event,
	})

	c.JSON(http.StatusOK, res)
}

func (e *eventController) UpdateEventHR(c *gin.Context) {
	eventID := c.Param("id")
	var event dto.CreaEventRequest

	if err := c.ShouldBindJSON(&event); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: "Payload type invalid",
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(event)
	if err != nil {
		errorMsg := helper.GetErrorMessage(err)
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: errorMsg,
		})
		return
	}

	userData, exists := c.Get("user")
	if !exists {
		errorhandler.HandleError(c, &errorhandler.UnauthorizedError{
			Message: "User not found",
		})
		return
	}

	userID, _ := userData.(map[string]interface{})["user_id"].(string)

	updatedEvent, err := e.serviceEvent.UpdateEventHR(&event, userID, eventID)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Event updated Successfully",
		Data:       updatedEvent,
	})

	c.JSON(http.StatusOK, res)
}

func (e *eventController) DeleteEventByID(c *gin.Context) {
	eventID := c.Param("id")
	err := e.serviceEvent.DeleteEvent(eventID)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Event deleted Successfully",
	})

	c.JSON(http.StatusOK, res)
}

func (e *eventController) GetEventsVendor(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		errorhandler.HandleError(c, &errorhandler.UnauthorizedError{
			Message: "User not found",
		})
		return
	}
	user_name, _ := userData.(map[string]interface{})["name"].(string)
	result, err := e.serviceEvent.GetAllEventsVendor(user_name)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Events retrieved Successfully",
		Data:       result,
	})

	c.JSON(http.StatusOK, res)
}

func (e *eventController) ConfirmDate(c *gin.Context) {
	var confirmeDate dto.ConfirmDateRequest
	eventID := c.Param("id")
	if err := c.ShouldBindJSON(&confirmeDate); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: "Payload type invalid",
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(confirmeDate)
	if err != nil {
		errorMsg := helper.GetErrorMessage(err)
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: errorMsg,
		})
		return
	}
	updatedEvent, err := e.serviceEvent.AcceptEventVendor(&confirmeDate, eventID)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Date confirmed Successfully",
		Data:       updatedEvent,
	})

	c.JSON(http.StatusOK, res)
}

func (e *eventController) RejectDates(c *gin.Context) {
	var rejectDate *dto.RejectRequest
	eventID := c.Param("id")
	if err := c.ShouldBindJSON(&rejectDate); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: "Payload type invalid",
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(rejectDate)
	if err != nil {
		errorMsg := helper.GetErrorMessage(err)
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: errorMsg,
		})
		return
	}
	updatedEvent, err := e.serviceEvent.RejectEventVendor(rejectDate, eventID)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Dates rejected Successfully",
		Data:       updatedEvent,
	})

	c.JSON(http.StatusOK, res)
}
