package models

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
