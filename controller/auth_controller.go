package controller

import (
	"booking-event-server/dto"
	errorhandler "booking-event-server/errorHandler"
	"booking-event-server/helper"
	"booking-event-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authController struct {
	service service.AuthService
}

func NewController(s service.AuthService) *authController {
	return &authController{
		service: s,
	}
}

func (a *authController) Register(c *gin.Context) {
	var user dto.RegisterRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: err.Error(),
		})
		return
	}

	if err := a.service.Register(&user); err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "User registered successfully",
	})

	c.JSON(http.StatusCreated, res)
}

func (a authController) Login(c *gin.Context) {
	var user dto.LoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: err.Error(),
		})
		return
	}

	token, err := a.service.Login(&user)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusAccepted,
		Message:    "User logged in successfully",
		Token:      &token,
	})

	c.JSON(http.StatusAccepted, res)
}
