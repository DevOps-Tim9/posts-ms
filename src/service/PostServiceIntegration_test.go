package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	"posts-ms/src/client"
	"posts-ms/src/dto/request"
	"posts-ms/src/entity"
	"posts-ms/src/rabbitmq"
	"posts-ms/src/repository"
	"posts-ms/src/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostServiceIntegrationTestSuite struct {
	suite.Suite
	service  PostService
	db       *gorm.DB
	likes    []entity.Like
	comments []entity.Comment
	posts    []entity.Post
}

func (suite *PostServiceIntegrationTestSuite) SetupSuite() {
	host := os.Getenv("DATABASE_DOMAIN")
	user := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	name := os.Getenv("DATABASE_SCHEMA")
	port := os.Getenv("DATABASE_PORT")

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		name,
		port,
	)

	db, _ := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	db.AutoMigrate(&entity.Like{Tbl: "likes"})
	db.AutoMigrate(&entity.Comment{Tbl: "comments"})
	db.AutoMigrate(&entity.Post{Tbl: "posts"})

	likeRepository := repository.LikeRepository{Database: db}
	commentRepository := repository.CommentRepository{Database: db}
	postRepository := repository.PostRepository{Database: db}

	mediaClient := client.MediaRESTClient{}

	suite.db = db

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	rabbit := rabbitmq.RMQProducer{
		ConnectionString: amqpServerURL,
	}

	channel, _ := rabbit.StartRabbitMQ()

	suite.service = PostService{
		LikeRepository:    likeRepository,
		CommentRepository: commentRepository,
		MediaClient:       mediaClient,
		RabbitMQChannel:   channel,
		PostRepository:    postRepository,
		Logger:            utils.Logger(),
	}

	suite.posts = []entity.Post{
		{
			Model: gorm.Model{
				ID:        10,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Description:  "Description",
			UserId:       8,
			TotalLikes:   1,
			TotalUnlikes: 1,
			ImageId:      1,
		},
		{
			Model: gorm.Model{
				ID:        11,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Description:  "Description",
			UserId:       789,
			TotalLikes:   1,
			TotalUnlikes: 1,
			ImageId:      1,
		},
	}

	suite.likes = []entity.Like{
		{
			Model: gorm.Model{
				ID:        10,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:   2,
			PostId:   10,
			LikeType: 1,
		},
		{
			Model: gorm.Model{
				ID:        20,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:   3,
			PostId:   10,
			LikeType: 2,
		},
	}

	suite.comments = []entity.Comment{
		{
			Model: gorm.Model{
				ID:        10,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:  2,
			PostId:  10,
			Content: "Comment1",
		},
		{
			Model: gorm.Model{
				ID:        20,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:  3,
			PostId:  10,
			Content: "Comment2",
		},
	}

	tx := suite.db.Begin()

	tx.Create(&suite.posts[0])
	tx.Create(&suite.likes[0])
	tx.Create(&suite.likes[1])
	tx.Create(&suite.comments[0])
	tx.Create(&suite.comments[1])

	tx.Commit()
}

func TestPostServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceIntegrationTestSuite))
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_GetById_PostDoesNotExist() {
	id := uint(9)

	post, err := suite.service.GetById(id, context.TODO())

	assert.Nil(suite.T(), post)
	assert.NotNil(suite.T(), err)
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_GetById_PostDoesExist() {
	id := uint(10)

	post, err := suite.service.GetById(id, context.TODO())

	assert.NotNil(suite.T(), post)
	assert.Equal(suite.T(), id, post.Id)
	assert.Equal(suite.T(), 1, post.TotalLikes)
	assert.Equal(suite.T(), 1, post.TotalUnlikes)
	assert.Nil(suite.T(), err)
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_GetPostById_PostDoesNotExist() {
	id := uint(9)

	_, err := suite.service.GetPostById(id, context.TODO())

	assert.NotNil(suite.T(), err)
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_GetPostById_PostDoesExist() {
	id := uint(10)

	post, err := suite.service.GetPostById(id, context.TODO())

	assert.NotNil(suite.T(), post)
	assert.Equal(suite.T(), id, post.ID)
	assert.Equal(suite.T(), 1, post.TotalLikes)
	assert.Equal(suite.T(), 1, post.TotalUnlikes)
	assert.Nil(suite.T(), err)
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_GetAllByUserId_PostDoesNotExist() {
	id := uint(99)

	posts := suite.service.GetAllByUserId(id, context.TODO())

	assert.Equal(suite.T(), 0, len(posts))
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_GetAllByUserId_PostDoesExist() {
	id := uint(8)

	posts := suite.service.GetAllByUserId(id, context.TODO())

	assert.Equal(suite.T(), 1, len(posts))
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_GetAllByUserIds_PostDoesNotExist() {
	ids := []uint{5, 6}

	posts := suite.service.GetAllByUserIds(ids, context.TODO())

	assert.Equal(suite.T(), 0, len(posts))
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_GetAllByUserIds_PostDoesExist() {
	ids := []uint{8, 6}

	posts := suite.service.GetAllByUserIds(ids, context.TODO())

	assert.Equal(suite.T(), 1, len(posts))
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_Delete_PostDoesNotExist() {
	id := uint(456)

	suite.service.Delete(id, context.TODO())

	assert.True(suite.T(), true)
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_Delete_PostDoesExist() {
	id := uint(11)
	userId := uint(789)

	suite.service.Delete(id, context.TODO())

	posts := suite.service.GetAllByUserId(userId, context.TODO())

	assert.Equal(suite.T(), 0, len(posts))
	assert.True(suite.T(), true)
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_Create_Successfully() {
	id := uint(21)

	newPost := request.PostDto{
		UserId:      id,
		Description: "Post",
	}

	file := []*multipart.FileHeader{
		{
			Filename: "a",
			Header:   textproto.MIMEHeader{},
			Size:     0,
		},
	}

	post, err := suite.service.Create(newPost, file, context.TODO())

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), post, err)
}

func (suite *PostServiceIntegrationTestSuite) TestIntegrationPostService_CreatePost_Successfully() {
	id := uint(21)

	newPost := entity.Post{
		UserId:      id,
		Description: "Post",
	}

	post, err := suite.service.CreatePost(newPost, context.TODO())

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), post, err)
	assert.Equal(suite.T(), "Post", post.Description)
	assert.Equal(suite.T(), id, post.UserId)
}
