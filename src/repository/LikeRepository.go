package repository

import (
	"posts-ms/src/entity"

	"gorm.io/gorm"
)

type ILikeRepository interface {
	Create(like entity.Like) (entity.Like, error)
	GetByUserIdAndPostId(uint, uint) (entity.Like, error)
	Delete(uint)
	DeleteByPostId(uint)
	GetAllByPostId(uint) []*entity.Like
}

type LikeRepository struct {
	Database *gorm.DB
}

func (r LikeRepository) GetAllByPostId(id uint) []*entity.Like {
	var likes = []*entity.Like{}

	r.Database.Find(&likes, "post_id = ?", id)

	return likes
}

func (r LikeRepository) Create(like entity.Like) (entity.Like, error) {
	error := r.Database.Save(&like).Error

	return like, error
}

func (r LikeRepository) GetByUserIdAndPostId(userId uint, postId uint) (entity.Like, error) {
	var like = entity.Like{}

	error := r.Database.
		Where("user_id = ?", userId).
		Where("post_id = ?", postId).
		First(&like).Error

	return like, error
}

func (r LikeRepository) Delete(id uint) {
	r.Database.Unscoped().Delete(&entity.Like{}, id)
}

func (r LikeRepository) DeleteByPostId(id uint) {
	r.Database.Unscoped().Where("post_id = ?", id).Delete(&entity.Like{})
}
