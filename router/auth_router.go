package router

import (
	"booking-event-server/config"
	"booking-event-server/controller"
	"booking-event-server/middleware"
	"booking-event-server/repository"
	"booking-event-server/service"

	"github.com/gin-gonic/gin"
)

func AuthRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewController(authService)

	auth := api.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
	api.GET("/user", middleware.Authentication, authController.UserDetail)
	api.GET("/vendors", middleware.Authentication, authController.AllVendors)
}
