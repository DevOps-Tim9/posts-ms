package repository

import (
	"errors"
	"posts-ms/src/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CommentRepositoryMock struct {
	mock.Mock
}

func (c CommentRepositoryMock) Create(comment entity.Comment) (entity.Comment, error) {
	comment.ID = 1

	return comment, nil
}

func (c CommentRepositoryMock) Delete(id uint) error {
	switch id {
	case 1:
		return nil
	case 2:
		return errors.New("")
	}

	return nil
}

func (c CommentRepositoryMock) DeleteByPostId(id uint) error {
	switch id {
	case 1:
		return nil
	case 2:
		return errors.New("")
	}

	return nil
}

func (c CommentRepositoryMock) GetAllByPostId(id uint) []*entity.Comment {
	switch id {
	case 1:
		return make([]*entity.Comment, 0)
	case 2:
		return []*entity.Comment{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Content: "Some text",
				UserId:  2,
				PostId:  1,
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Content: "Some text 2",
				UserId:  5,
				PostId:  2,
			}}

	}

	return nil
}
