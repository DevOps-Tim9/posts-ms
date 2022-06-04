package request

type PostDto struct {
	Description string `json:"description" validate:"required"`
	UserId      uint   `json:"userId" validate:"required"`
}
