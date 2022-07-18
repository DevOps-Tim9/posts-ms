package entity

import (
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Description  string `gorm:"default:null"`
	ImageId      uint
	UserId       uint
	TotalLikes   int
	TotalUnlikes int
	Likes        []Like
	Comments     []Comment
	Tbl          string `gorm:"-"`
}

func CreatePost(dto request.PostDto) Post {
	return Post{
		Description:  dto.Description,
		UserId:       dto.UserId,
		TotalLikes:   0,
		TotalUnlikes: 0,
		ImageId:      0,
	}
}

func (post Post) CreateDto() *response.PostDto {
	return &response.PostDto{
		Id:           post.ID,
		Description:  post.Description,
		UserId:       post.UserId,
		ImageId:      post.ImageId,
		TotalLikes:   post.TotalLikes,
		TotalUnlikes: post.TotalUnlikes,
		Likes:        transformLikesToDtos(post.Likes),
		Comments:     transformCommentsToDtos(post.Comments),
	}
}

func transformLikesToDtos(likes []Like) []response.LikeDto {
	var likesDto = []response.LikeDto{}

	for _, value := range likes {
		like := value.CreateDto()

		likesDto = append(likesDto, *like)
	}

	return likesDto
}

func transformCommentsToDtos(likes []Comment) []response.CommentDto {
	var commentsDto = []response.CommentDto{}

	for _, value := range likes {
		comment := value.CreateDto()

		commentsDto = append(commentsDto, *comment)
	}

	return commentsDto
}

func (post *Post) SetImageId(imageId uint) {
	post.ImageId = imageId
}
