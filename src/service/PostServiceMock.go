package service

import (
	"context"
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

func (p PostServiceMock) Create(request.PostDto, []*multipart.FileHeader, context.Context) (*response.PostDto, error) {
	return nil, nil
}

func (p PostServiceMock) CreatePost(entity.Post, context.Context) (*entity.Post, error) {
	return nil, nil
}

func (p PostServiceMock) Delete(uint, context.Context) {

}

func (p PostServiceMock) GetById(id uint, ctx context.Context) (*response.PostDto, error) {
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

func (p PostServiceMock) GetPostById(id uint, ctx context.Context) (*entity.Post, error) {
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

func (p PostServiceMock) GetAllByUserId(uint, context.Context) []*response.PostDto {
	return nil
}

func (p PostServiceMock) GetAllByUserIds([]uint, context.Context) []*response.PostDto {
	return nil
}
