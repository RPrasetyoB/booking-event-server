package middleware

import (
	errorhandler "booking-event-server/errorHandler"

	"github.com/gin-gonic/gin"
)

func authorization(allowedRoles []int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, exists := c.Get("user")
		if !exists {
			errorhandler.HandleError(c, &errorhandler.AccessForbiddenError{
				Message: "User data not found in context",
			})
			c.Abort()
			return
		}

		roleID := userData.(map[string]interface{})["role_id"]
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if roleID == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			errorhandler.HandleError(c, &errorhandler.AccessForbiddenError{
				Message: "Access forbidden: Role not allowed",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

var HrAuth = authorization([]int64{1})
var VendorAuth = authorization([]int64{2})
