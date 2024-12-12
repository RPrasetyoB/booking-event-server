package errorhandler

import (
	"booking-event-server/dto"
	"booking-event-server/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	var statusCode int

	switch err.(type) {
	case *NotFoundError:
		statusCode = http.StatusNotFound
	case *BadRequestError:
		statusCode = http.StatusBadRequest
	case *InternalServerError:
		statusCode = http.StatusInternalServerError
	case *UnauthorizedError:
		statusCode = http.StatusUnauthorized
	case *AccessForbiddenError:
		statusCode = http.StatusForbidden
	}

	response := helper.Response(dto.ResponseParams{
		StatusCode: statusCode,
		Message:    err.Error(),
	})

	c.JSON(statusCode, response)
}
