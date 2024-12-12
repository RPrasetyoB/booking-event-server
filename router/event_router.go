package router

import (
	"booking-event-server/config"
	"booking-event-server/controller"
	"booking-event-server/middleware"
	"booking-event-server/repository"
	"booking-event-server/service"

	"github.com/gin-gonic/gin"
)

func EventRouter(api *gin.RouterGroup) {
	eventRepository := repository.NewEventRepository(config.DB)
	dateRepository := repository.NewDatesRepository(config.DB)
	eventService := service.NewEventService(eventRepository, dateRepository)
	eventController := controller.NewEventController(eventService)

	event := api.Group("/event")
	event.POST("/", middleware.Authentication, eventController.CreateEvent)
}
