package response

type CommentDto struct {
	Id      uint   `json:"id"`
	PostId  uint   `json:"postId" validate:"required"`
	UserId  uint   `json:"userId" validate:"required"`
	Content string `json:"content" validate:"required"`
}
