package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoAuthHeader         = fmt.Errorf("authorization header is missing")
	ErrTokenExpired         = fmt.Errorf("token has expired")
	ErrInvalidSigningMethod = fmt.Errorf("unexpected signing method")
)

func GetToken(c *gin.Context) (map[string]interface{}, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, ErrNoAuthHeader
	}

	tokenString := strings.Split(authHeader, " ")[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		secretKey := os.Getenv("JWT_SECRET")

		return []byte(secretKey), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, jwt.ErrTokenExpired
		}
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errors.New("invalid user_id in token")
		}
		userID = fmt.Sprintf("%d", int64(userIDFloat))
	}

	user_name, ok := claims["name"].(string)
	if !ok {
		userIDFloat, ok := claims["name"].(float64)
		if !ok {
			return nil, errors.New("invalid name in token")
		}
		userID = fmt.Sprintf("%d", int64(userIDFloat))
	}

	roleIDFloat, ok := claims["role_id"].(float64)
	if !ok {
		return nil, errors.New("invalid role_id in token")
	}
	roleID := int64(roleIDFloat)

	result := map[string]interface{}{
		"user_id": userID,
		"role_id": roleID,
		"name":    user_name,
	}
	return result, nil
}
