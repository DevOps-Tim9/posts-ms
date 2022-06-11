package service

import (
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"
	"posts-ms/src/entity"
	"posts-ms/src/repository"

	"github.com/sirupsen/logrus"
)

type ILikeService interface {
	Create(request.LikeDto) (*response.LikeDto, error)
	Delete(uint, uint)
	GetAllByPostId(uint) []*response.LikeDto
}

type LikeService struct {
	LikeRepository repository.ILikeRepository
	PostService    IPostService
	Logger         *logrus.Entry
}

func (s LikeService) GetAllByPostId(id uint) []*response.LikeDto {
	s.Logger.Info("Getting likes for post")

	likes := s.LikeRepository.GetAllByPostId(id)

	return s.transformListOfDAOToListOfDTO(likes)
}

func (s LikeService) Create(dto request.LikeDto) (*response.LikeDto, error) {
	s.Logger.Info("Creating like")

	like, error := s.LikeRepository.GetByUserIdAndPostId(dto.UserId, dto.PostId)

	if error == nil {
		like.LikeType = entity.TypeOfLike(dto.LikeType)
	} else {
		like = entity.CreateLike(dto)
	}

	newLike, _ := s.LikeRepository.Create(like)

	post, error := s.PostService.GetPostById(dto.PostId)

	if error != nil {
		return nil, error
	}

	totalPositive := 0
	totalNegative := 0

	for _, item := range post.Likes {
		if item.LikeType == 1 {
			totalPositive = totalPositive + 1
		} else {
			totalNegative = totalNegative + 1
		}
	}

	post.TotalLikes = totalPositive
	post.TotalUnlikes = totalNegative

	s.PostService.CreatePost(*post)

	return newLike.CreateDto(), error
}

func (s LikeService) Delete(userId uint, postId uint) {
	s.Logger.Info("Deleting like")

	like, error := s.LikeRepository.GetByUserIdAndPostId(userId, postId)

	if error != nil {
		return
	}

	s.LikeRepository.Delete(like.ID)

	post, error := s.PostService.GetPostById(postId)

	if error != nil {
		return
	}

	if like.LikeType == 1 {
		post.TotalLikes = post.TotalLikes - 1
	} else {
		post.TotalUnlikes = post.TotalUnlikes - 1
	}

	s.PostService.CreatePost(*post)
}

func (s LikeService) transformListOfDAOToListOfDTO(likes []*entity.Like) []*response.LikeDto {
	var likesDto = []*response.LikeDto{}

	for _, value := range likes {
		likeDto := value.CreateDto()

		likesDto = append(likesDto, likeDto)
	}

	return likesDto
}
