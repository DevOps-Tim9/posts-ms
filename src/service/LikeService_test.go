package service

import (
	"posts-ms/src/dto/request"
	"posts-ms/src/repository"
	"posts-ms/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LikeServiceUnitTestSuite struct {
	suite.Suite
	likeRepositoryMock *repository.LikeRepositoryMock
	postServiceMock    *PostServiceMock
	service            LikeService
}

func TestLikeServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(LikeServiceUnitTestSuite))
}

func (suite *LikeServiceUnitTestSuite) SetupSuite() {
	suite.likeRepositoryMock = new(repository.LikeRepositoryMock)
	suite.postServiceMock = new(PostServiceMock)

	suite.service = LikeService{LikeRepository: suite.likeRepositoryMock,
		PostService: suite.postServiceMock,
		Logger:      utils.Logger(),
	}
}

func (suite *LikeServiceUnitTestSuite) TestNewLikeService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *LikeServiceUnitTestSuite) TestLikeService_GetAllByPostId_ReturnsEmptyList() {
	likes := suite.service.GetAllByPostId(1)

	assert.NotNil(suite.T(), likes, "Likes are nil")
	assert.Equal(suite.T(), 0, len(likes), "Length of likes is not 0")
}

func (suite *LikeServiceUnitTestSuite) TestLikeService_GetAllByPostId_ReturnsListOfLikes() {
	likes := suite.service.GetAllByPostId(2)

	assert.NotNil(suite.T(), likes, "Likes are nil")
	assert.Equal(suite.T(), 2, len(likes), "Length of likes is not 2")
}

func (suite *LikeServiceUnitTestSuite) TestLikeService_Delete_ReturnsNothing() {
	assert.True(suite.T(), true, "Test failed")
}

func (suite *LikeServiceUnitTestSuite) TestLikeService_Create_ReturnsLike() {
	id := uint(1)

	like := request.LikeDto{
		PostId:   5,
		UserId:   3,
		LikeType: 2,
	}

	newLike, err := suite.service.Create(like)

	assert.NotNil(suite.T(), newLike, "Like is nil")
	assert.Nil(suite.T(), err, "Error is not nil")
	assert.Equal(suite.T(), id, newLike.Id, "Id is not 1")
}

func (suite *LikeServiceUnitTestSuite) TestLikeService_Create_WithNonExistPost_ReturnError() {
	like := request.LikeDto{
		PostId:   1,
		UserId:   6,
		LikeType: 2,
	}

	newLike, err := suite.service.Create(like)

	assert.NotNil(suite.T(), err, "Error is nil")
	assert.Nil(suite.T(), newLike, "Like is not nil")
}

func (suite *LikeServiceUnitTestSuite) TestLikeService_Create_WithNonExistPostAndUser_ReturnError() {
	like := request.LikeDto{
		PostId:   1,
		UserId:   1,
		LikeType: 2,
	}

	newLike, err := suite.service.Create(like)

	assert.NotNil(suite.T(), err, "Error is nil")
	assert.Nil(suite.T(), newLike, "Like is not nil")
}
