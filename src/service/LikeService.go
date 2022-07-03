package service

import (
	"context"
	"fmt"
	"posts-ms/src/client"
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"
	"posts-ms/src/entity"
	"posts-ms/src/rabbitmq"
	"posts-ms/src/repository"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type ILikeService interface {
	Create(request.LikeDto, context.Context) (*response.LikeDto, error)
	Delete(uint, uint, context.Context)
	GetAllByPostId(uint, context.Context) []*response.LikeDto
}

type LikeService struct {
	LikeRepository  repository.ILikeRepository
	PostService     IPostService
	Logger          *logrus.Entry
	UserRESTClient  client.IUserRESTClient
	RabbitMQChannel *amqp.Channel
}

func (s LikeService) GetAllByPostId(id uint, ctx context.Context) []*response.LikeDto {
	s.Logger.Info("Getting likes for post")

	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get all likes for specific post")

	defer span.Finish()

	likes := s.LikeRepository.GetAllByPostId(id, ctx)

	return s.transformListOfDAOToListOfDTO(likes)
}

func (s LikeService) Create(dto request.LikeDto, ctx context.Context) (*response.LikeDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Create new like for specific post")

	defer span.Finish()

	s.Logger.Info("Creating like")

	like, error := s.LikeRepository.GetByUserIdAndPostId(dto.UserId, dto.PostId, ctx)

	if error == nil {
		like.LikeType = entity.TypeOfLike(dto.LikeType)
	} else {
		like = entity.CreateLike(dto)
	}

	newLike, _ := s.LikeRepository.Create(like, ctx)

	post, error := s.PostService.GetPostById(dto.PostId, ctx)

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

	s.PostService.CreatePost(*post, ctx)

	s.AddNotification(int(dto.UserId), int(post.UserId), dto.LikeType, ctx)

	return newLike.CreateDto(), error
}

func (s LikeService) Delete(userId uint, postId uint, ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Delete like for specific post from specific user")

	defer span.Finish()

	s.Logger.Info("Deleting like")

	like, error := s.LikeRepository.GetByUserIdAndPostId(userId, postId, ctx)

	if error != nil {
		return
	}

	s.LikeRepository.Delete(like.ID, ctx)

	post, error := s.PostService.GetPostById(postId, ctx)

	if error != nil {
		return
	}

	if like.LikeType == 1 {
		post.TotalLikes = post.TotalLikes - 1
	} else {
		post.TotalUnlikes = post.TotalUnlikes - 1
	}

	s.PostService.CreatePost(*post, ctx)
}

func (s LikeService) transformListOfDAOToListOfDTO(likes []*entity.Like) []*response.LikeDto {
	var likesDto = []*response.LikeDto{}

	for _, value := range likes {
		likeDto := value.CreateDto()

		likesDto = append(likesDto, likeDto)
	}

	return likesDto
}

func (s LikeService) AddNotification(fromId int, toId int, likeType int, ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Notify user about new post")

	defer span.Finish()

	userFrom, _ := s.UserRESTClient.GetUser(fromId, ctx)
	userTo, _ := s.UserRESTClient.GetUser(toId, ctx)

	var notification request.NotificationDTO
	messageType := request.Like
	if likeType == 1 {
		notification = request.NotificationDTO{Message: fmt.Sprintf("%s liked your post.", userFrom.Username), UserAuth0ID: userTo.Auth0ID, NotificationType: &messageType}
	} else {
		notification = request.NotificationDTO{Message: fmt.Sprintf("%s disliked your post.", userFrom.Username), UserAuth0ID: userTo.Auth0ID, NotificationType: &messageType}
	}

	rabbitmq.AddNotification(&notification, s.RabbitMQChannel, ctx)
}
