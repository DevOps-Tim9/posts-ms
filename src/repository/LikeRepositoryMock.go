package repository

import (
	"context"
	"errors"
	"posts-ms/src/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type LikeRepositoryMock struct {
	mock.Mock
}

func (l LikeRepositoryMock) Create(like entity.Like, ctx context.Context) (entity.Like, error) {
	like.ID = 1

	return like, nil
}

func (l LikeRepositoryMock) GetByUserIdAndPostId(userId uint, postId uint, ctx context.Context) (entity.Like, error) {
	if userId == 1 && postId == 1 {
		return entity.Like{}, errors.New("")
	} else {
		return entity.Like{
			Model: gorm.Model{
				ID: 1,
			},
			UserId:   2,
			PostId:   1,
			LikeType: 1,
		}, nil
	}
}

func (l LikeRepositoryMock) Delete(uint, context.Context) {
}

func (l LikeRepositoryMock) DeleteByPostId(uint, context.Context) {
}

func (l LikeRepositoryMock) GetAllByPostId(id uint, ctx context.Context) []*entity.Like {
	switch id {
	case 1:
		return make([]*entity.Like, 0)
	case 2:
		return []*entity.Like{
			{
				Model: gorm.Model{
					ID: 1,
				},
				UserId:   2,
				PostId:   1,
				LikeType: 1,
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				UserId:   3,
				PostId:   1,
				LikeType: 2,
			}}

	}

	return nil
}
