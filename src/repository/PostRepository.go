package repository

import (
	"posts-ms/src/entity"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IPostRepository interface {
	Create(post entity.Post) (entity.Post, error)
	Delete(uint)
	GetById(uint) (*entity.Post, error)
	GetAllByUserId(uint) []*entity.Post
	GetAllByUserIds([]uint) []*entity.Post
}

type PostRepository struct {
	Database *gorm.DB
}

func (r PostRepository) GetById(id uint) (*entity.Post, error) {
	var post = entity.Post{}

	error := r.Database.Preload("Likes").First(&post, id).Error

	return &post, error
}

func (r PostRepository) GetAllByUserId(id uint) []*entity.Post {
	var posts = []*entity.Post{}

	r.Database.Preload("Likes").Preload("Comments").Order("created_at desc").Find(&posts, "user_id = ?", id)

	return posts
}

func (r PostRepository) GetAllByUserIds(ids []uint) []*entity.Post {
	var posts = []*entity.Post{}

	r.Database.Preload("Likes").Preload("Comments").Order("created_at desc").Find(&posts, "user_id = any(?)", pq.Array(ids))

	return posts
}

func (r PostRepository) Create(post entity.Post) (entity.Post, error) {
	error := r.Database.Save(&post).Error

	return post, error
}

func (r PostRepository) Delete(id uint) {
	r.Database.Unscoped().Select(clause.Associations).Delete(&entity.Post{}, id)
}
