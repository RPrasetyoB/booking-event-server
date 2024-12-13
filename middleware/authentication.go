package middleware

import (
	errorhandler "booking-event-server/errorHandler"
	"booking-event-server/helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication(c *gin.Context) {
	userData, err := helper.GetToken(c)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			errorhandler.HandleError(c, &errorhandler.UnauthorizedError{
				Message: "Token expired, please re-login",
			})
			c.Abort()
			return
		}

		errorhandler.HandleError(c, &errorhandler.UnauthorizedError{
			Message: "Unauthorized",
		})
		c.Abort()
		return
	}

	c.Set("user", userData)
	c.Next()
}
