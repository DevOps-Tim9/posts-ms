package repository

import (
	"context"
	"posts-ms/src/entity"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

type ILikeRepository interface {
	Create(entity.Like, context.Context) (entity.Like, error)
	GetByUserIdAndPostId(uint, uint, context.Context) (entity.Like, error)
	Delete(uint, context.Context)
	DeleteByPostId(uint, context.Context)
	GetAllByPostId(uint, context.Context) []*entity.Like
}

type LikeRepository struct {
	Database *gorm.DB
}

func (r LikeRepository) GetAllByPostId(id uint, ctx context.Context) []*entity.Like {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Get all likes for specific post")

	defer span.Finish()

	var likes = []*entity.Like{}

	r.Database.Find(&likes, "post_id = ?", id)

	return likes
}

func (r LikeRepository) Create(like entity.Like, ctx context.Context) (entity.Like, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Create new like")

	defer span.Finish()

	error := r.Database.Save(&like).Error

	return like, error
}

func (r LikeRepository) GetByUserIdAndPostId(userId uint, postId uint, ctx context.Context) (entity.Like, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Get like likes for specific post from specific user")

	defer span.Finish()

	var like = entity.Like{}

	error := r.Database.
		Where("user_id = ?", userId).
		Where("post_id = ?", postId).
		First(&like).Error

	return like, error
}

func (r LikeRepository) Delete(id uint, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository -Delete by id")

	defer span.Finish()

	r.Database.Unscoped().Delete(&entity.Like{}, id)
}

func (r LikeRepository) DeleteByPostId(id uint, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Delete by post id")

	defer span.Finish()

	r.Database.Unscoped().Where("post_id = ?", id).Delete(&entity.Like{})
}
