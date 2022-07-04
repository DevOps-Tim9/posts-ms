package service

import (
	"context"
	"mime/multipart"
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

type IPostService interface {
	Create(request.PostDto, []*multipart.FileHeader, context.Context) (*response.PostDto, error)
	CreatePost(entity.Post, context.Context) (*entity.Post, error)
	Delete(uint, context.Context)
	GetById(uint, context.Context) (*response.PostDto, error)
	GetPostById(uint, context.Context) (*entity.Post, error)
	GetAllByUserId(uint, context.Context) []*response.PostDto
	GetAllByUserIds([]uint, context.Context) []*response.PostDto
}

type PostService struct {
	PostRepository    repository.IPostRepository
	LikeRepository    repository.ILikeRepository
	CommentRepository repository.ICommentRepository
	MediaClient       client.IMediaClient
	RabbitMQChannel   *amqp.Channel
	Logger            *logrus.Entry
}

func (s PostService) GetById(id uint, ctx context.Context) (*response.PostDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get post by id")

	defer span.Finish()

	s.Logger.Info("Getting post by id")

	post, err := s.PostRepository.GetById(id, ctx)

	if err != nil {
		return nil, err
	}

	return post.CreateDto(), err
}

func (s PostService) GetPostById(id uint, ctx context.Context) (*entity.Post, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get post by id")

	defer span.Finish()

	s.Logger.Info("Getting post by id")

	return s.PostRepository.GetById(id, ctx)
}

func (s PostService) GetAllByUserId(id uint, ctx context.Context) []*response.PostDto {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get all posts by user id")

	defer span.Finish()

	s.Logger.Info("Getting posts by user")

	posts := s.PostRepository.GetAllByUserId(id, ctx)

	return s.transformListOfDAOToListOfDTO(posts)
}

func (s PostService) GetAllByUserIds(ids []uint, ctx context.Context) []*response.PostDto {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Get posts by user ids")

	defer span.Finish()

	s.Logger.Info("Getting posts by users")

	posts := s.PostRepository.GetAllByUserIds(ids, ctx)

	return s.transformListOfDAOToListOfDTO(posts)
}

func (s PostService) Create(dto request.PostDto, images []*multipart.FileHeader, ctx context.Context) (*response.PostDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Create post")

	defer span.Finish()

	s.Logger.Info("Creating post")

	post := entity.CreatePost(dto)

	file, _ := images[0].Open()

	s.Logger.Info("Sending request on media-ms for creating media")
	imageId, err := s.MediaClient.Upload(file, ctx)

	if err != nil {
		return nil, err
	}

	post.SetImageId(imageId)

	newPost, err := s.PostRepository.Create(post, ctx)

	return newPost.CreateDto(), err
}

func (s PostService) CreatePost(post entity.Post, ctx context.Context) (*entity.Post, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Create post")

	defer span.Finish()

	post, err := s.PostRepository.Create(post, ctx)

	return &post, err
}

func (s PostService) Delete(id uint, ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Service - Delete post by id")

	defer span.Finish()

	s.Logger.Info("Deleting post")

	post, error := s.PostRepository.GetById(id, ctx)

	if error != nil {
		return
	}

	s.Logger.Info("Sending request on media-ms for deleting media")
	rabbitmq.DeleteImage(post.ImageId, s.RabbitMQChannel, ctx)

	s.Logger.Info("Deleting likes for post")
	s.LikeRepository.DeleteByPostId(id, ctx)

	s.Logger.Info("Deleting comments for post")
	s.CommentRepository.DeleteByPostId(id, ctx)

	s.PostRepository.Delete(id, ctx)
}

func (s PostService) transformListOfDAOToListOfDTO(posts []*entity.Post) []*response.PostDto {
	var postsDto = []*response.PostDto{}

	for _, value := range posts {
		postDto := value.CreateDto()

		postsDto = append(postsDto, postDto)
	}

	return postsDto
}
