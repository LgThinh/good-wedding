package mapper

import (
	"github.com/google/uuid"
	"good-wedding/pkg/model"
)

func MapUser(username string) *model.User {
	return &model.User{
		UserName: username,
	}
}

func MapComment(userID uuid.UUID, req *model.CommentRequest) *model.Comment {
	return &model.Comment{
		ObjectID: *req.ObjectID,
		UserID:   userID,
		Comment:  *req.Comment,
	}
}

func MapWeddingWish(userID uuid.UUID, req *model.WeddingWishRequest) *model.WeddingWish {
	return &model.WeddingWish{
		UserID:  userID,
		Comment: *req.Comment,
	}
}
