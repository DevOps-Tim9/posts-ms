package service

import (
	"posts-ms/src/dto/request"
	"posts-ms/src/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommentServiceUnitTestSuite struct {
	suite.Suite
	commentRepositoryMock *repository.CommentRepositoryMock
	service               CommentService
}

func TestCommentServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(CommentServiceUnitTestSuite))
}

func (suite *CommentServiceUnitTestSuite) SetupSuite() {
	suite.commentRepositoryMock = new(repository.CommentRepositoryMock)

	suite.service = CommentService{CommentRepository: suite.commentRepositoryMock}
}

func (suite *CommentServiceUnitTestSuite) TestNewCommentService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *CommentServiceUnitTestSuite) TestCommentService_GetAllByPostId_ReturnsEmptyList() {
	comments := suite.service.GetAllByPostId(1)

	assert.NotNil(suite.T(), comments, "Comments are nil")
	assert.Equal(suite.T(), 0, len(comments), "Length of comments is not 0")
}

func (suite *CommentServiceUnitTestSuite) TestCommentService_GetAllByPostId_ReturnsListOfComments() {
	comments := suite.service.GetAllByPostId(2)

	assert.NotNil(suite.T(), comments, "Comments are nil")
	assert.Equal(suite.T(), 2, len(comments), "Length of comments is not 2")
}

func (suite *CommentServiceUnitTestSuite) TestCommentService_Delete_CommentNotExist() {
	suite.service.Delete(2)

	assert.True(suite.T(), true, "")
}

func (suite *CommentServiceUnitTestSuite) TestCommentService_Create_ReturnsComment() {
	id := uint(1)

	comment := request.CommentDto{
		Content: "Some text",
		UserId:  1,
	}

	newComment, err := suite.service.Create(comment)

	assert.NotNil(suite.T(), newComment, "Comment is nil")
	assert.Equal(suite.T(), id, newComment.Id, "Comment id is not 1")
	assert.Nil(suite.T(), err, "Error is not nil")
}
