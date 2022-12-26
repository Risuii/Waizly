package models

import "time"

type Account struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Email     string    `json:"email" validate:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}
