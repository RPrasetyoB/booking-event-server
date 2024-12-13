package service

import (
	"booking-event-server/dto"
	"booking-event-server/entity"
	errorhandler "booking-event-server/errorHandler"
	"booking-event-server/helper"
	"booking-event-server/repository"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (string, error)
	UserProfile(userID string) (*dto.UserResponse, error)
	GetAllVendors() ([]*dto.UserResponse, error)
}

type authService struct {
	repository repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) *authService {
	return &authService{
		repository: r,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	nameExist, _ := s.repository.FindName(req.Name)

	if nameExist != nil {
		return &errorhandler.BadRequestError{Message: "name already taken"}
	}

	passwordHash, err := helper.HashPassword(req.Password)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	var roleID int
	switch req.Role {
	case "hr":
		roleID = 1
	case "vendor":
		roleID = 2
	}

	nanoid, _ := helper.GenerateNanoId()

	user := entity.User{
		ID:       nanoid,
		Name:     req.Name,
		Password: passwordHash,
		Role_id:  roleID,
	}

	if err := s.repository.Register(&user); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}

func (s *authService) Login(req *dto.LoginRequest) (string, error) {
	user, err := s.repository.FindName(req.Name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", &errorhandler.BadRequestError{Message: "name or password invalid"}
		}
		return "", &errorhandler.InternalServerError{Message: err.Error()}
	}
	if !helper.ComparePassword(req.Password, user.Password) {
		return "", &errorhandler.BadRequestError{Message: "name or password invalid"}
	}

	token, err := helper.GenerateToken(*user)
	if err != nil {
		return "", &errorhandler.InternalServerError{Message: err.Error()}
	}

	return token, nil
}

func (s *authService) UserProfile(userID string) (*dto.UserResponse, error) {
	user, err := s.repository.FindUserById(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errorhandler.BadRequestError{Message: "user not found"}
		}
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	response := &dto.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Role_id: user.Role_id,
	}

	return response, nil
}

func (s *authService) GetAllVendors() ([]*dto.UserResponse, error) {
	users, err := s.repository.GetAllVendors()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errorhandler.BadRequestError{Message: "users not found"}
		}
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	var responses []*dto.UserResponse
	for _, user := range users {
		responses = append(responses, &dto.UserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Role_id: user.Role_id,
		})
	}

	return responses, nil
}
