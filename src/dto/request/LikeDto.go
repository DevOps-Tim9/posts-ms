package request

type LikeDto struct {
	PostId   uint `json:"postId" validate:"required"`
	UserId   uint `json:"userId" validate:"required"`
	LikeType int  `json:"likeType" validate:"required,min=1,max=2"`
}
