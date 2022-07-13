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

type LikeServiceIntegrationTestSuite struct {
	suite.Suite
	service LikeService
	db      *gorm.DB
	likes   []entity.Like
	posts   []entity.Post
}

func (suite *LikeServiceIntegrationTestSuite) SetupSuite() {
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
	postRepository := repository.PostRepository{Database: db}

	userRESTClient := client.UserRESTClient{}

	suite.db = db

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	rabbit := rabbitmq.RMQProducer{
		ConnectionString: amqpServerURL,
	}

	channel, _ := rabbit.StartRabbitMQ()

	suite.service = LikeService{
		PostService:     PostService{PostRepository: postRepository, Logger: utils.Logger()},
		UserRESTClient:  userRESTClient,
		RabbitMQChannel: channel,
		LikeRepository:  likeRepository,
		Logger:          utils.Logger(),
	}

	suite.posts = []entity.Post{
		{
			Model: gorm.Model{
				ID:        1000,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Description:  "Description",
			UserId:       1,
			TotalLikes:   1,
			TotalUnlikes: 1,
			ImageId:      1,
		},
		{
			Model: gorm.Model{
				ID:        2000,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Description:  "Description",
			UserId:       1,
			TotalLikes:   1,
			TotalUnlikes: 1,
			ImageId:      2,
		},
	}
	suite.likes = []entity.Like{
		{
			Model: gorm.Model{
				ID:        1000,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:   2,
			PostId:   1000,
			LikeType: 1,
		},
		{
			Model: gorm.Model{
				ID:        2000,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:   3,
			PostId:   1000,
			LikeType: 2,
		},
		{
			Model: gorm.Model{
				ID:        3000,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:   2,
			PostId:   2000,
			LikeType: 1,
		},
		{
			Model: gorm.Model{
				ID:        4000,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			UserId:   3,
			PostId:   2000,
			LikeType: 2,
		},
	}

	tx := suite.db.Begin()

	tx.Create(&suite.posts[0])
	tx.Create(&suite.posts[1])
	tx.Create(&suite.likes[0])
	tx.Create(&suite.likes[1])
	tx.Create(&suite.likes[2])
	tx.Create(&suite.likes[3])

	tx.Commit()
}

func TestLikeServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(LikeServiceIntegrationTestSuite))
}

func (suite *LikeServiceIntegrationTestSuite) TestIntegrationLikeService_GetAllByPostId_PostDoesNotExist() {
	id := uint(10)

	likes := suite.service.GetAllByPostId(id, context.TODO())

	assert.Equal(suite.T(), 0, len(likes))
}

func (suite *LikeServiceIntegrationTestSuite) TestIntegrationLikeService_GetAllByPostId_PostDoesExist() {
	id := uint(1000)

	likes := suite.service.GetAllByPostId(id, context.TODO())

	assert.GreaterOrEqual(suite.T(), len(likes), 1)
}

func (suite *LikeServiceIntegrationTestSuite) TestIntegrationLikeService_Delete_LikeDoesNotExist() {
	postId := uint(9999)
	existPostId := uint(1000)
	userId := uint(10)

	suite.service.Delete(userId, postId, context.TODO())

	post, err := suite.service.PostService.GetById(existPostId, context.TODO())

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, post.TotalLikes)
	assert.Equal(suite.T(), 1, post.TotalUnlikes)
}

func (suite *LikeServiceIntegrationTestSuite) TestIntegrationLikeService_Delete_LikeDoesExist() {
	postId := uint(2000)
	userId := uint(2)

	suite.service.Delete(userId, postId, context.TODO())

	post, err := suite.service.PostService.GetById(postId, context.TODO())

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, post.TotalLikes)
	assert.Equal(suite.T(), 1, post.TotalUnlikes)
}

func (suite *LikeServiceIntegrationTestSuite) TestIntegrationLikeService_Create_Successfully() {
	id := uint(1000)

	likeDto := request.LikeDto{
		PostId:   id,
		UserId:   1,
		LikeType: 1,
	}

	like, err := suite.service.Create(likeDto, context.TODO())

	likes := suite.service.GetAllByPostId(id, context.TODO())

	post, _ := suite.service.PostService.GetById(id, context.TODO())

	assert.Equal(suite.T(), 2, post.TotalLikes)
	assert.Equal(suite.T(), 1, post.TotalUnlikes)
	assert.Equal(suite.T(), 3, len(likes))
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), like)
	assert.Equal(suite.T(), 1, like.LikeType)

	suite.service.Delete(like.UserId, like.PostId, context.TODO())
}
