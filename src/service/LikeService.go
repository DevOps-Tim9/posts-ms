package service

import (
	"fmt"
	"posts-ms/src/client"
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"
	"posts-ms/src/entity"
	"posts-ms/src/rabbitmq"
	"posts-ms/src/repository"

	"github.com/streadway/amqp"
)

type ILikeService interface {
	Create(request.LikeDto) (*response.LikeDto, error)
	Delete(uint, uint)
	GetAllByPostId(uint) []*response.LikeDto
}

type LikeService struct {
	LikeRepository  repository.ILikeRepository
	PostService     IPostService
	UserRESTClient  client.IUserRESTClient
	RabbitMQChannel *amqp.Channel
}

func (s LikeService) GetAllByPostId(id uint) []*response.LikeDto {
	likes := s.LikeRepository.GetAllByPostId(id)

	return s.transformListOfDAOToListOfDTO(likes)
}

func (s LikeService) Create(dto request.LikeDto) (*response.LikeDto, error) {
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

	s.AddNotification(int(dto.UserId), int(post.UserId), dto.LikeType)

	return newLike.CreateDto(), error
}

func (s LikeService) Delete(userId uint, postId uint) {

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

func (s LikeService) AddNotification(fromId int, toId int, likeType int) {
	userFrom, _ := s.UserRESTClient.GetUser(fromId)
	userTo, _ := s.UserRESTClient.GetUser(toId)

	var notification request.NotificationDTO
	messageType := request.Like
	if likeType == 1 {
		notification = request.NotificationDTO{Message: fmt.Sprintf("%s liked your post.", userFrom.Username), UserAuth0ID: userTo.Auth0ID, NotificationType: &messageType}
	} else {
		notification = request.NotificationDTO{Message: fmt.Sprintf("%s disliked your post.", userFrom.Username), UserAuth0ID: userTo.Auth0ID, NotificationType: &messageType}
	}

	rabbitmq.AddNotification(&notification, s.RabbitMQChannel)
}
