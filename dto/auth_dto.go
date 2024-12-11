package dto

type RegisterRequest struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
