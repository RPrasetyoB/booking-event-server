package repository

import (
	"booking-event-server/entity"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserById(userId string) (*entity.User, error)
	FindName(name string) (*entity.User, error)
	Register(req *entity.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) FindUserById(userId string) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "user_id = ?", userId).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindName(name string) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "name = ?", name).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) Register(user *entity.User) error {
	err := r.db.Create(&user).Error

	return err
}
