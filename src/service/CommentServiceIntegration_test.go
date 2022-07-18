package service

import (
	"context"
	"fmt"
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

type CommentServiceIntegrationTestSuite struct {
	suite.Suite
	service  CommentService
	db       *gorm.DB
	comments []entity.Comment
	posts    []entity.Post
}

func (suite *CommentServiceIntegrationTestSuite) SetupSuite() {
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

	db.AutoMigrate(&entity.Comment{Tbl: "comments"})
	db.AutoMigrate(&entity.Like{Tbl: "likes"})
	db.AutoMigrate(&entity.Post{Tbl: "posts"})

	commentRepository := repository.CommentRepository{Database: db}
	postrepository := repository.PostRepository{Database: db}

	postService := PostService{PostRepository: postrepository, Logger: utils.Logger()}

	userRESTClient := client.UserRESTClient{}

	suite.db = db

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	rabbit := rabbitmq.RMQProducer{
		ConnectionString: amqpServerURL,
	}

	channel, _ := rabbit.StartRabbitMQ()

	suite.service = CommentService{
		PostService:       postService,
		UserRESTClient:    userRESTClient,
		RabbitMQChannel:   channel,
		CommentRepository: commentRepository,
		Logger:            utils.Logger(),
	}

	suite.posts = []entity.Post{
		{
			Model: gorm.Model{
				ID:        100,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Description:  "Description",
			UserId:       1,
			TotalLikes:   0,
			TotalUnlikes: 0,
			ImageId:      1,
		},
		{
			Model: gorm.Model{
				ID:        200,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Description:  "Description",
			UserId:       1,
			TotalLikes:   0,
			TotalUnlikes: 0,
			ImageId:      1,
		},
	}
	suite.comments = []entity.Comment{
		{
			Model: gorm.Model{
				ID:        100,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Content: "Comment 1",
			UserId:  2,
			PostId:  100,
		},
		{
			Model: gorm.Model{
				ID:        200,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Content: "Comment 2",
			UserId:  3,
			PostId:  100,
		},
		{
			Model: gorm.Model{
				ID:        300,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Content: "Comment 3",
			UserId:  3,
			PostId:  200,
		},
	}

	tx := suite.db.Begin()

	tx.Create(&suite.posts[0])
	tx.Create(&suite.posts[1])
	tx.Create(&suite.comments[0])
	tx.Create(&suite.comments[1])
	tx.Create(&suite.comments[2])

	tx.Commit()
}

func TestCommentServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CommentServiceIntegrationTestSuite))
}

func (suite *CommentServiceIntegrationTestSuite) TestIntegrationCommentService_GetAllByPostId_PostDoesNotExist() {
	id := uint(10)

	comments := suite.service.GetAllByPostId(id, context.TODO())

	assert.Equal(suite.T(), 0, len(comments))
}

func (suite *CommentServiceIntegrationTestSuite) TestIntegrationCommentService_GetAllByPostId_PostDoesExist() {
	id := uint(100)

	comments := suite.service.GetAllByPostId(id, context.TODO())

	assert.GreaterOrEqual(suite.T(), len(comments), 2)
}

func (suite *CommentServiceIntegrationTestSuite) TestIntegrationCommentService_Delete_CommentDoesNotExist() {
	id := uint(10)

	suite.service.Delete(id, context.TODO())

	assert.True(suite.T(), true)
}

func (suite *CommentServiceIntegrationTestSuite) TestIntegrationCommentService_Delete_CommentDoesExist() {
	commentId := uint(300)
	postId := uint(200)

	suite.service.Delete(commentId, context.TODO())

	comments := suite.service.GetAllByPostId(postId, context.TODO())

	assert.Equal(suite.T(), 0, len(comments))
}

func (suite *CommentServiceIntegrationTestSuite) TestIntegrationCommentService_Create_Successfully() {
	id := uint(100)

	commentDto := request.CommentDto{
		PostId:  id,
		UserId:  2,
		Content: "Comment",
	}

	comment, err := suite.service.Create(commentDto, context.TODO())

	comments := suite.service.GetAllByPostId(id, context.TODO())

	assert.Equal(suite.T(), len(comments), 3)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), comment)
	assert.Equal(suite.T(), "Comment", comment.Content)

	suite.service.Delete(comment.Id, context.TODO())
}
