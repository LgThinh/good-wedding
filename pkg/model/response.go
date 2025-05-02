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
