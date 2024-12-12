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
	service service.EventService
}

func NewEventController(s service.EventService) *eventController {
	return &eventController{
		service: s,
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

	newEvent, err := e.service.CreateEvent(&event, userID)
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
