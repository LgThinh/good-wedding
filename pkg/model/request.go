package model

import (
	"github.com/google/uuid"
	"good-wedding/pkg/model/paging"
)

type CommentRequest struct {
	ObjectID *uuid.UUID `json:"object_id" binding:"required"`
	UserName *string    `json:"user_name" binding:"required"`
	Comment  *string    `json:"comment" binding:"required"`
}

type WeddingWishRequest struct {
	UserName *string `json:"user_name" binding:"required"`
	Comment  *string `json:"comment" binding:"required"`
}

type CommentFilterRequest struct {
	FromDate *int64  `json:"from_date" form:"from_date"`
	ToDate   *int64  `json:"to_date" form:"to_date"`
	ObjectID *string `json:"object_id" form:"object_id"`
}

type CommentFilter struct {
	CommentFilterRequest
	Pager *paging.Pager
}

type WeddingWishFilterRequest struct {
	FromDate *int64 `json:"from_date" form:"from_date"`
	ToDate   *int64 `json:"to_date" form:"to_date"`
}

type WeddingWishFilter struct {
	WeddingWishFilterRequest
	Pager *paging.Pager
}

type UserFilterRequest struct {
	FromDate *int64  `json:"from_date" form:"from_date"`
	ToDate   *int64  `json:"to_date" form:"to_date"`
	UserName *string `json:"user_name" form:"user_name"`
}

type UserFilter struct {
	UserFilterRequest
	Pager *paging.Pager
}

type ObjectMediaFilterRequest struct {
	FromDate   *int64  `json:"from_date" form:"from_date"`
	ToDate     *int64  `json:"to_date" form:"to_date"`
	Name       *string `json:"name" form:"name"`
	Url        *string `json:"url" form:"url"`
	ObjectType *string `json:"object_type" form:"object_type"`
	FileType   *string `json:"file_type" form:"file_type"`
}

type ObjectMediaFilter struct {
	ObjectMediaFilterRequest
	Pager *paging.Pager
}
