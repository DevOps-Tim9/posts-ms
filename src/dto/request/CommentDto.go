package request

type CommentDto struct {
	PostId  uint   `json:"postId" validate:"required"`
	UserId  uint   `json:"userId" validate:"required"`
	Content string `json:"content" validate:"required"`
}
