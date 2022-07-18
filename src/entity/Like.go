package entity

import (
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"

	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserId   uint `gorm:"not null;default:null"`
	PostId   uint `gorm:"not null;default:null"`
	Post     Post
	LikeType TypeOfLike `gorm:"not null;default:null"`

	Tbl string `gorm:"-"`
}

func CreateLike(dto request.LikeDto) Like {
	return Like{
		UserId:   dto.UserId,
		PostId:   dto.PostId,
		LikeType: TypeOfLike(dto.LikeType),
	}
}

func (like Like) CreateDto() *response.LikeDto {
	return &response.LikeDto{
		Id:       like.ID,
		PostId:   like.PostId,
		UserId:   like.UserId,
		LikeType: int(like.LikeType),
	}
}
