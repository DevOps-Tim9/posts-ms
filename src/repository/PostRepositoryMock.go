package repository

import (
	"errors"
	"posts-ms/src/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type PostRepositoryMock struct {
	mock.Mock
}

func (p PostRepositoryMock) Create(post entity.Post) (entity.Post, error) {
	post.ID = 1

	return post, nil
}

func (p PostRepositoryMock) Delete(uint) {

}

func (p PostRepositoryMock) GetById(id uint) (*entity.Post, error) {
	if id == 1 {
		return nil, errors.New("")
	} else {
		return &entity.Post{
			Model: gorm.Model{
				ID: 2,
			},
			UserId:       2,
			Description:  "Some text",
			ImageId:      1,
			TotalLikes:   0,
			TotalUnlikes: 0,
		}, nil
	}
}

func (p PostRepositoryMock) GetAllByUserId(id uint) []*entity.Post {
	if id == 1 {
		return []*entity.Post{}
	} else {
		return []*entity.Post{
			{
				Model: gorm.Model{
					ID: 1,
				},
				UserId:       2,
				Description:  "Some text",
				ImageId:      1,
				TotalLikes:   0,
				TotalUnlikes: 0,
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				UserId:       2,
				Description:  "Some text",
				ImageId:      2,
				TotalLikes:   0,
				TotalUnlikes: 0,
			},
		}
	}
}

func (p PostRepositoryMock) GetAllByUserIds(ids []uint) []*entity.Post {
	if ids[0] == 1 && ids[1] == 2 {
		return []*entity.Post{}
	} else {
		return []*entity.Post{
			{
				Model: gorm.Model{
					ID: 1,
				},
				UserId:       2,
				Description:  "Some text",
				ImageId:      1,
				TotalLikes:   0,
				TotalUnlikes: 0,
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				UserId:       2,
				Description:  "Some text",
				ImageId:      2,
				TotalLikes:   0,
				TotalUnlikes: 0,
			},
		}
	}
}
