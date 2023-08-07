package models

import "time"

type Posts struct {
	Id          uint          `json:"id"`
	Description string        `json:"description"`
	Image       string        `json:"image"`
	LikesCount  int           `json:"like" gorm:"default:0"`
	UserId      uint          `json:"user_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Comments    []PostCommets `json:"comments" gorm:"foreignKey:PostId"`
	Likes       []PostLikes   `json:"likes" gorm:"foreignKey:PostId"`
}
