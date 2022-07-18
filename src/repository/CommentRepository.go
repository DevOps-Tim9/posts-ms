package repository

import (
	"context"
	"posts-ms/src/entity"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	Create(entity.Comment, context.Context) (entity.Comment, error)
	Delete(uint, context.Context) error
	DeleteByPostId(uint, context.Context) error
	GetAllByPostId(uint, context.Context) []*entity.Comment
}

type CommentRepository struct {
	Database *gorm.DB
}

func (r CommentRepository) GetAllByPostId(id uint, ctx context.Context) []*entity.Comment {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Get all comments for specific post")

	defer span.Finish()

	var comments = []*entity.Comment{}

	r.Database.Find(&comments, "post_id = ?", id)

	return comments
}

func (r CommentRepository) Create(comment entity.Comment, ctx context.Context) (entity.Comment, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Create new comment for specific post")

	defer span.Finish()

	error := r.Database.Save(&comment).Error

	return comment, error
}

func (r CommentRepository) Delete(id uint, ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Delete comment by id")

	defer span.Finish()

	r.Database.Unscoped().Delete(&entity.Comment{}, id)

	return nil
}

func (r CommentRepository) DeleteByPostId(id uint, ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Delete comment by post id")

	defer span.Finish()

	r.Database.Unscoped().Where("post_id = ?", id).Delete(&entity.Comment{})

	return nil
}
