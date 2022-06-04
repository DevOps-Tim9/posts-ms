package response

type LikeDto struct {
	Id       uint `json:"id"`
	PostId   uint `json:"postId" validate:"required"`
	UserId   uint `json:"userId" validate:"required"`
	LikeType int  `json:"likeType" validate:"required"`
}
