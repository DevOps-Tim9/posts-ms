package service

import (
	"mime/multipart"
	"net/textproto"
	"posts-ms/src/client"
	"posts-ms/src/dto/request"
	"posts-ms/src/entity"
	"posts-ms/src/repository"
	"posts-ms/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PostServiceUnitTestSuite struct {
	suite.Suite
	postRepositoryMock  *repository.PostRepositoryMock
	mediaRestClientMock *client.MediaRestClientMock
	service             PostService
}

func TestPostServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceUnitTestSuite))
}

func (suite *PostServiceUnitTestSuite) SetupSuite() {
	suite.postRepositoryMock = new(repository.PostRepositoryMock)
	suite.mediaRestClientMock = new(client.MediaRestClientMock)

	suite.service = PostService{PostRepository: suite.postRepositoryMock,
		MediaClient: suite.mediaRestClientMock,
		Logger:      utils.Logger(),
	}
}

func (suite *PostServiceUnitTestSuite) TestNewPostService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *PostServiceUnitTestSuite) TestPostService_GetById_ReturnPost() {
	id := uint(2)

	post, err := suite.service.GetById(2)

	assert.NotNil(suite.T(), post, "Post is nil")
	assert.Nil(suite.T(), err, "Error is not nil")
	assert.Equal(suite.T(), id, post.Id, "Post id is not 2")
}

func (suite *PostServiceUnitTestSuite) TestPostService_GetById_ReturnError() {
	post, err := suite.service.GetById(1)

	assert.NotNil(suite.T(), err, "Error is nil")
	assert.Nil(suite.T(), post, "Post is not nil")
}

func (suite *PostServiceUnitTestSuite) TestPostService_GetAllByUserId_ReturnEmptyList() {
	posts := suite.service.GetAllByUserId(1)

	assert.NotNil(suite.T(), posts, "Posts are nil")
	assert.Equal(suite.T(), 0, len(posts), "Length of posts not 0")
}

func (suite *PostServiceUnitTestSuite) TestPostService_GetAllByUserId_ReturnListOfPosts() {
	posts := suite.service.GetAllByUserId(2)

	assert.NotNil(suite.T(), posts, "Posts are nil")
	assert.Equal(suite.T(), 2, len(posts), "Length of posts not 2")
}

func (suite *PostServiceUnitTestSuite) TestPostService_GetAllByUsersId_ReturnEmptyList() {
	posts := suite.service.GetAllByUserIds([]uint{1, 2})

	assert.NotNil(suite.T(), posts, "Posts are nil")
	assert.Equal(suite.T(), 0, len(posts), "Length of posts not 0")
}

func (suite *PostServiceUnitTestSuite) TestPostService_GetAllByUsersId_ReturnListOfPosts() {
	posts := suite.service.GetAllByUserIds([]uint{2, 6})

	assert.NotNil(suite.T(), posts, "Posts are nil")
	assert.Equal(suite.T(), 2, len(posts), "Length of posts not 2")
}

func (suite *PostServiceUnitTestSuite) TestPostService_CreatePost_ReturnPost() {
	id := uint(1)

	post := entity.Post{
		Description:  "Some text",
		ImageId:      1,
		UserId:       1,
		TotalLikes:   0,
		TotalUnlikes: 0,
	}

	newPost, err := suite.service.CreatePost(post)

	assert.NotNil(suite.T(), newPost, "Posts are nil")
	assert.Nil(suite.T(), err, "Error is not nil")
	assert.Equal(suite.T(), id, newPost.ID, "Post id is not 2")
}

func (suite *PostServiceUnitTestSuite) TestPostService_Create_ReturnPost() {
	id := uint(1)

	post := request.PostDto{
		Description: "Some text",
		UserId:      1,
	}

	newPost, err := suite.service.Create(post, []*multipart.FileHeader{
		{
			Filename: "",
			Header:   textproto.MIMEHeader{},
			Size:     0,
		},
	})

	assert.NotNil(suite.T(), newPost, "Posts are nil")
	assert.Nil(suite.T(), err, "Error is not nil")
	assert.Equal(suite.T(), id, newPost.Id, "Post id is not 2")
}

func (suite *PostServiceUnitTestSuite) TestPostService_TransformListOfDAOToListOfDTO_ReturnEmptyList() {
	posts := suite.service.transformListOfDAOToListOfDTO([]*entity.Post{})

	assert.NotNil(suite.T(), posts, "Posts are nil")
	assert.Equal(suite.T(), 0, len(posts), "Length of posts not 0")
}

func (suite *PostServiceUnitTestSuite) TestPostService_TransformListOfDAOToListOfDTO_ReturnListOfPosts() {
	posts := suite.service.transformListOfDAOToListOfDTO([]*entity.Post{{
		Description:  "Some text",
		ImageId:      1,
		UserId:       1,
		TotalLikes:   0,
		TotalUnlikes: 0,
	}})

	assert.NotNil(suite.T(), posts, "Posts are nil")
	assert.Equal(suite.T(), 1, len(posts), "Length of posts not 1")
}
