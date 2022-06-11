package service

import (
	"mime/multipart"
	"posts-ms/src/client"
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"
	"posts-ms/src/entity"
	"posts-ms/src/rabbitmq"
	"posts-ms/src/repository"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type IPostService interface {
	Create(request.PostDto, []*multipart.FileHeader) (*response.PostDto, error)
	CreatePost(entity.Post) (*entity.Post, error)
	Delete(uint)
	GetById(uint) (*response.PostDto, error)
	GetPostById(uint) (*entity.Post, error)
	GetAllByUserId(uint) []*response.PostDto
	GetAllByUserIds([]uint) []*response.PostDto
}

type PostService struct {
	PostRepository    repository.IPostRepository
	LikeRepository    repository.ILikeRepository
	CommentRepository repository.ICommentRepository
	MediaClient       client.IMediaClient
	RabbitMQChannel   *amqp.Channel
	Logger            *logrus.Entry
}

func (s PostService) GetById(id uint) (*response.PostDto, error) {
	s.Logger.Info("Getting post by id")

	post, err := s.PostRepository.GetById(id)

	if err != nil {
		return nil, err
	}

	return post.CreateDto(), err
}

func (s PostService) GetPostById(id uint) (*entity.Post, error) {
	s.Logger.Info("Getting post by id")

	return s.PostRepository.GetById(id)
}

func (s PostService) GetAllByUserId(id uint) []*response.PostDto {
	s.Logger.Info("Getting posts by user")

	posts := s.PostRepository.GetAllByUserId(id)

	return s.transformListOfDAOToListOfDTO(posts)
}

func (s PostService) GetAllByUserIds(ids []uint) []*response.PostDto {
	s.Logger.Info("Getting posts by users")

	posts := s.PostRepository.GetAllByUserIds(ids)

	return s.transformListOfDAOToListOfDTO(posts)
}

func (s PostService) Create(dto request.PostDto, images []*multipart.FileHeader) (*response.PostDto, error) {
	s.Logger.Info("Creating post")

	post := entity.CreatePost(dto)

	file, _ := images[0].Open()

	s.Logger.Info("Sending request on media-ms for creating media")
	imageId, err := s.MediaClient.Upload(file)

	if err != nil {
		return nil, err
	}

	post.SetImageId(imageId)

	newPost, err := s.PostRepository.Create(post)

	return newPost.CreateDto(), err
}

func (s PostService) CreatePost(post entity.Post) (*entity.Post, error) {
	post, err := s.PostRepository.Create(post)

	return &post, err
}

func (s PostService) Delete(id uint) {
	s.Logger.Info("Deleting post")

	post, error := s.PostRepository.GetById(id)

	if error != nil {
		return
	}

	s.Logger.Info("Sending request on media-ms for deleting media")
	rabbitmq.DeleteImage(post.ImageId, s.RabbitMQChannel)

	s.Logger.Info("Deleting likes for post")
	s.LikeRepository.DeleteByPostId(id)

	s.Logger.Info("Deleting comments for post")
	s.CommentRepository.DeleteByPostId(id)

	s.PostRepository.Delete(id)
}

func (s PostService) transformListOfDAOToListOfDTO(posts []*entity.Post) []*response.PostDto {
	var postsDto = []*response.PostDto{}

	for _, value := range posts {
		postDto := value.CreateDto()

		postsDto = append(postsDto, postDto)
	}

	return postsDto
}
