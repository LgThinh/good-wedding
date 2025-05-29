package model

import "github.com/google/uuid"

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
	Filter  *CommentFilter         `json:"filter"`
	Records []*CommentDataResponse `json:"data"`
}

type CommentDataResponse struct {
	InitTime  string     `json:"init_time"`
	ObjectID  *uuid.UUID `json:"object_id"`
	ObjectUrl string     `json:"object_url"`
	UserID    uuid.UUID  `json:"user_id"`
	UserName  string     `json:"user_name"`
	Comment   string     `json:"comment"`
}
type WeddingWishFilterResult struct {
	Filter  *WeddingWishFilter         `json:"filter"`
	Records []*WeddingWishDataResponse `json:"data"`
}

type WeddingWishDataResponse struct {
	InitTime string    `json:"init_time"`
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
	Comment  string    `json:"comment"`
}

type UserFilterResult struct {
	Filter  *UserFilter `json:"filter"`
	Records []*User     `json:"data"`
}

type ObjectMediaFilterResult struct {
	Filter  *ObjectMediaFilter `json:"filter"`
	Records []*ObjectMedia     `json:"data"`
}

type GetOneImageResult struct {
	Meta *MetaData   `json:"meta"`
	Data ObjectMedia `json:"data"`
}

type GetOneVideoResult struct {
	Meta *MetaData   `json:"meta"`
	Data ObjectMedia `json:"data"`
}
