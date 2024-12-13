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
	authRepository := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(authRepository)
	eventService := service.NewEventService(eventRepository, dateRepository, authRepository)
	eventController := controller.NewEventController(eventService, authService)

	event := api.Group("/event")
	event.POST("/hr", middleware.Authentication, eventController.CreateEvent)
	event.GET("/hr", middleware.Authentication, middleware.HrAuth, eventController.GetEventsHRbyUserID)
	event.GET("/:id", middleware.Authentication, eventController.GetEventbyID)
	event.PUT("/hr/:id", middleware.Authentication, middleware.HrAuth, eventController.UpdateEventHR)
	event.DELETE("/:id", middleware.Authentication, eventController.DeleteEventByID)
	event.GET("/vendor", middleware.Authentication, middleware.VendorAuth, eventController.GetEventsVendor)
	event.PATCH("/vendor/:id", middleware.Authentication, middleware.VendorAuth, eventController.ConfirmDate)
	event.PATCH("/vendor/reject/:id", middleware.Authentication, middleware.VendorAuth, eventController.RejectDates)
}
