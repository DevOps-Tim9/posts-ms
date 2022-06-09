package repository

import (
	"posts-ms/src/entity"

	"gorm.io/gorm"
)

type ICommentRepository interface {
	Create(comment entity.Comment) (entity.Comment, error)
	Delete(uint) error
	DeleteByPostId(uint) error
	GetAllByPostId(uint) []*entity.Comment
}

type CommentRepository struct {
	Database *gorm.DB
}

func (r CommentRepository) GetAllByPostId(id uint) []*entity.Comment {
	var comments = []*entity.Comment{}

	r.Database.Find(&comments, "post_id = ?", id)

	return comments
}

func (r CommentRepository) Create(comment entity.Comment) (entity.Comment, error) {
	error := r.Database.Save(&comment).Error

	return comment, error
}

func (r CommentRepository) Delete(id uint) error {
	r.Database.Unscoped().Delete(&entity.Comment{}, id)

	return nil
}

func (r CommentRepository) DeleteByPostId(id uint) error {
	r.Database.Unscoped().Where("post_id = ?", id).Delete(&entity.Comment{})

	return nil
}
