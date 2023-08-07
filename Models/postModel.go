package models

import "time"

type Posts struct {
	Id        uint      `json:"id"`
	Title     string    `json:"username"`
	Body      string    `json:"password"`
	Image     string    `json:"image"`
	UserId    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
