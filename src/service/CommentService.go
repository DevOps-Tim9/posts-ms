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

type ICommentService interface {
	Create(request.CommentDto) (*response.CommentDto, error)
	Delete(uint)
	GetAllByPostId(uint) []*response.CommentDto
}

type CommentService struct {
	CommentRepository repository.ICommentRepository
	PostService       IPostService
	UserRESTClient    client.IUserRESTClient
	RabbitMQChannel   *amqp.Channel
}

func (s CommentService) GetAllByPostId(id uint) []*response.CommentDto {
	comments := s.CommentRepository.GetAllByPostId(id)

	return transformListOfDAOToListOfDTO(comments)
}

func (s CommentService) Create(dto request.CommentDto) (*response.CommentDto, error) {
	comment := entity.CreateComment(dto)
	post, error := s.PostService.GetPostById(dto.PostId)

	newComment, error := s.CommentRepository.Create(comment)

	s.AddNotification(int(dto.UserId), int(post.UserId))

	return newComment.CreateDto(), error
}

func (s CommentService) Delete(id uint) {
	s.CommentRepository.Delete(id)
}

func transformListOfDAOToListOfDTO(comments []*entity.Comment) []*response.CommentDto {
	var commentsDto = []*response.CommentDto{}

	for _, value := range comments {
		commentDto := value.CreateDto()

		commentsDto = append(commentsDto, commentDto)
	}

	return commentsDto
}

func (s CommentService) AddNotification(fromId int, toId int) {
	userFrom, _ := s.UserRESTClient.GetUser(fromId)
	userTo, _ := s.UserRESTClient.GetUser(toId)

	messageType := request.Comment
	notification := request.NotificationDTO{Message: fmt.Sprintf("%s commented on your post.", userFrom.Username), UserAuth0ID: userTo.Auth0ID, NotificationType: &messageType}

	rabbitmq.AddNotification(&notification, s.RabbitMQChannel)
}
