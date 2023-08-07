package models

type PostLikes struct {
	Id     uint `json:"id"`
	PostId uint `json:"post_id"`
	UserId uint `json:"user_id"`
}
