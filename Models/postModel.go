package models

import "time"

type Posts struct {
	Id          uint          `json:"id"`
	Description string        `json:"description" validate:"required,max=255"`
	Image       string        `json:"image"`
	LikesCount  int           `json:"like" gorm:"default:0"`
	UserId      uint          `json:"user_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Comments    []PostCommets `json:"comments" gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
	Likes       []PostLikes   `json:"likes" gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
}
