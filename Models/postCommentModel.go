package models

import "time"

type PostCommets struct {
	Id        uint      `json:"id"`
	Comment   string    `json:"comment"`
	PostId    uint      `json:"post_id"`
	UserId    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
