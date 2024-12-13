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
			Message: "Payload type invalid",
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		errorMsg := helper.GetErrorMessage(err)
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: errorMsg,
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
			Message: "Payload type invalid",
		})
		return
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		errorMsg := helper.GetErrorMessage(err)
		errorhandler.HandleError(c, &errorhandler.BadRequestError{
			Message: errorMsg,
		})
		return
	}

	token, err := a.service.Login(&user)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "User logged in successfully",
		Token:      &token,
	})

	c.JSON(http.StatusOK, res)
}

func (a authController) UserDetail(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		errorhandler.HandleError(c, &errorhandler.UnauthorizedError{
			Message: "User not found",
		})
		return
	}

	userID, _ := userData.(map[string]interface{})["user_id"].(string)
	userData, err := a.service.UserProfile(userID)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "User Data retrieved successfully",
		Data:       userData,
	})

	c.JSON(http.StatusOK, res)
}

func (a authController) AllVendors(c *gin.Context) {
	userData, err := a.service.GetAllVendors()
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Vendors retrieved successfully",
		Data:       userData,
	})

	c.JSON(http.StatusOK, res)
}
