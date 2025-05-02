package model

type UploadImageSuccessResponse struct {
	Meta *MetaData      `json:"meta"`
	Data UploadImageUrl `json:"data"`
}
type UploadImageUrl struct {
	Url string `json:"url"`
}

type UploadVideoSuccessResponse struct {
	Meta *MetaData      `json:"meta"`
	Data UploadVideoUrl `json:"data"`
}
type UploadVideoUrl struct {
	Url string `json:"url"`
}

type StringResponse struct {
	Meta *MetaData `json:"meta"`
	Data string    `json:"data"`
}

type CommentFilterResult struct {
	Filter  *CommentFilter `json:"filter"`
	Records []*Comment     `json:"data"`
}

type WeddingWishFilterResult struct {
	Filter  *WeddingWishFilter `json:"filter"`
	Records []*WeddingWish     `json:"data"`
}

type UserFilterResult struct {
	Filter  *UserFilter `json:"filter"`
	Records []*User     `json:"data"`
}

type ObjectMediaFilterResult struct {
	Filter  *ObjectMediaFilter `json:"filter"`
	Records []*ObjectMedia     `json:"data"`
}
