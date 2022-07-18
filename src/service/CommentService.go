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

type ICommentService interface {
	Create(request.CommentDto, context.Context) (*response.CommentDto, error)
	Delete(uint, context.Context)
	GetAllByPostId(uint, context.Context) []*response.CommentDto
}

type CommentService struct {
	CommentRepository repository.ICommentRepository
	Logger            *logrus.Entry
	PostService       IPostService
	UserRESTClient    client.IUserRESTClient
	RabbitMQChannel   *amqp.Channel
}

func (s CommentService) GetAllByPostId(id uint, ctx context.Context) []*response.CommentDto {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get all comments for specific post")

	defer span.Finish()

	s.Logger.Info("Getting comments for post")

	comments := s.CommentRepository.GetAllByPostId(id, ctx)

	return transformListOfDAOToListOfDTO(comments)
}

func (s CommentService) Create(dto request.CommentDto, ctx context.Context) (*response.CommentDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Create new comment for specific post")

	defer span.Finish()

	s.Logger.Info("Creating comment")

	comment := entity.CreateComment(dto)

	post, _ := s.PostService.GetPostById(dto.PostId, ctx)

	newComment, err := s.CommentRepository.Create(comment, ctx)

	s.AddNotification(int(dto.UserId), int(post.UserId), ctx)

	return newComment.CreateDto(), err
}

func (s CommentService) Delete(id uint, ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Delete comment by id")

	defer span.Finish()

	s.Logger.Info("Deleting comment")

	s.CommentRepository.Delete(id, ctx)
}

func transformListOfDAOToListOfDTO(comments []*entity.Comment) []*response.CommentDto {
	var commentsDto = []*response.CommentDto{}

	for _, value := range comments {
		commentDto := value.CreateDto()
		commentsDto = append(commentsDto, commentDto)
	}

	return commentsDto
}

func (s CommentService) AddNotification(fromId int, toId int, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Service - Notify user about new comment")

	defer span.Finish()

	userFrom, _ := s.UserRESTClient.GetUser(fromId, ctx)
	userTo, _ := s.UserRESTClient.GetUser(toId, ctx)

	messageType := request.Comment
	notification := request.NotificationDTO{Message: fmt.Sprintf("%s commented on your post.", userFrom.Username), UserAuth0ID: userTo.Auth0ID, NotificationType: &messageType}

	rabbitmq.AddNotification(&notification, s.RabbitMQChannel, ctx)
}
