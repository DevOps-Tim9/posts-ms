package response

type PostDto struct {
	Id           uint         `json:"id"`
	Description  string       `json:"description" validate:"required"`
	UserId       uint         `json:"userId" validate:"required"`
	ImageId      uint         `json:"imageId" validate:"required"`
	TotalLikes   int          `json:"totalLikes" validate:"required"`
	TotalUnlikes int          `json:"totalUnlikes" validate:"required"`
	Likes        []LikeDto    `json:"likes"`
	Comments     []CommentDto `json:"comments"`
}
