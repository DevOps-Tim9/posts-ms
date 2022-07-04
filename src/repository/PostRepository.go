package repository

import (
	"context"
	"posts-ms/src/entity"

	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IPostRepository interface {
	Create(entity.Post, context.Context) (entity.Post, error)
	Delete(uint, context.Context)
	GetById(uint, context.Context) (*entity.Post, error)
	GetAllByUserId(uint, context.Context) []*entity.Post
	GetAllByUserIds([]uint, context.Context) []*entity.Post
}

type PostRepository struct {
	Database *gorm.DB
}

func (r PostRepository) GetById(id uint, ctx context.Context) (*entity.Post, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Get post by id")

	defer span.Finish()

	var post = entity.Post{}

	error := r.Database.Preload("Likes").First(&post, id).Error

	return &post, error
}

func (r PostRepository) GetAllByUserId(id uint, ctx context.Context) []*entity.Post {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Get all posts by user id")

	defer span.Finish()

	var posts = []*entity.Post{}

	r.Database.Preload("Likes").Preload("Comments").Order("created_at desc").Find(&posts, "user_id = ?", id)

	return posts
}

func (r PostRepository) GetAllByUserIds(ids []uint, ctx context.Context) []*entity.Post {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Get all posts by user ids")

	defer span.Finish()

	var posts = []*entity.Post{}

	r.Database.Preload("Likes").Preload("Comments").Order("created_at desc").Find(&posts, "user_id = any(?)", pq.Array(ids))

	return posts
}

func (r PostRepository) Create(post entity.Post, ctx context.Context) (entity.Post, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Create post")

	defer span.Finish()

	error := r.Database.Save(&post).Error

	return post, error
}

func (r PostRepository) Delete(id uint, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Repository - Delete post by id")

	defer span.Finish()

	r.Database.Unscoped().Select(clause.Associations).Delete(&entity.Post{}, id)
}
