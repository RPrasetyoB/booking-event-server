package helper

import (
	"booking-event-server/entity"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(passwordHash), err
}

func ComparePassword(password string, hashedPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password)) == nil
}

func GenerateToken(user entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"role_id": user.Role_id,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}
	return tokenString, nil
}
