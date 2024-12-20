package dto

type RegisterRequest struct {
	Name     string `validate:"required,min=2" json:"name"`
	Role     string `validate:"required,oneof=hr vendor" json:"role"`
	Password string `validate:"required,min=6" json:"password"`
}

type LoginRequest struct {
	Name     string `validate:"required" json:"name"`
	Password string `validate:"required" json:"password"`
}

type UserResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Role_id int    `json:"role_id"`
}
