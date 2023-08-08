package models

import (
	"time"
)

type Users struct {
	Id        uint      `json:"id"`
	Username  string    `json:"username" validate:"required,min=3,usernameValid"`
	Password  string    `json:"password" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
