package request

type SearchPostPageableDto struct {
	Ids []uint `json:"userIds" validate:"required"`
}
