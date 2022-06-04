package service

import (
	"errors"
	"mime/multipart"
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"
	"posts-ms/src/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type PostServiceMock struct {
	mock.Mock
}

func (p PostServiceMock) Create(request.PostDto, []*multipart.FileHeader) (*response.PostDto, error) {
	return nil, nil
}

func (p PostServiceMock) CreatePost(entity.Post) (*entity.Post, error) {
	return nil, nil
}

func (p PostServiceMock) Delete(uint) {

}

func (p PostServiceMock) GetById(id uint) (*response.PostDto, error) {
	switch id {
	case 1:
		return nil, errors.New("")
	}
	return &response.PostDto{
		Id:           1,
		Description:  "Some text",
		UserId:       1,
		ImageId:      1,
		TotalLikes:   2,
		TotalUnlikes: 2,
	}, nil
}

func (p PostServiceMock) GetPostById(id uint) (*entity.Post, error) {
	switch id {
	case 1:
		return nil, errors.New("")
	}
	return &entity.Post{
		Model: gorm.Model{
			ID: 1,
		},
		Description:  "Some text",
		UserId:       1,
		ImageId:      1,
		TotalLikes:   2,
		TotalUnlikes: 2,
	}, nil
}

func (p PostServiceMock) GetAllByUserId(uint) []*response.PostDto {
	return nil
}

func (p PostServiceMock) GetAllByUserIds([]uint) []*response.PostDto {
	return nil
}
