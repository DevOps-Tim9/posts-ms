package service

import (
	"posts-ms/src/dto/request"
	"posts-ms/src/dto/response"
	"posts-ms/src/entity"
	"posts-ms/src/repository"
)

type ICommentService interface {
	Create(request.CommentDto) (*response.CommentDto, error)
	Delete(uint)
	GetAllByPostId(uint) []*response.CommentDto
}

type CommentService struct {
	CommentRepository repository.ICommentRepository
}

func (s CommentService) GetAllByPostId(id uint) []*response.CommentDto {
	comments := s.CommentRepository.GetAllByPostId(id)

	return transformListOfDAOToListOfDTO(comments)
}

func (s CommentService) Create(dto request.CommentDto) (*response.CommentDto, error) {
	comment := entity.CreateComment(dto)

	newComment, error := s.CommentRepository.Create(comment)

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
