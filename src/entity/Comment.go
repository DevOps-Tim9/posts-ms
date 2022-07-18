package entity

import (
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"not null;default:null"`
	UserId  uint   `gorm:"not null;default:null"`
	PostId  uint   `gorm:"not null;default:null"`
	Post    Post

	Tbl string `gorm:"-"`
}

func CreateComment(dto request.CommentDto) Comment {
	return Comment{
		UserId:  dto.UserId,
		PostId:  dto.PostId,
		Content: dto.Content,
	}
}

func (comment Comment) CreateDto() *response.CommentDto {
	return &response.CommentDto{
		Id:      comment.ID,
		PostId:  comment.PostId,
		UserId:  comment.UserId,
		Content: comment.Content,
	}
}
