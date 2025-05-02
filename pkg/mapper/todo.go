package mapper

import (
	"github.com/google/uuid"
	"good-wedding/pkg/model"
	"good-wedding/pkg/utils"
)

func MapTodo(req model.CreateTodoRequest, adminID uuid.UUID) *model.Todo {
	return &model.Todo{
		BaseModel: model.BaseModel{
			CreatorID: &adminID,
		},
		Name:        *req.Name,
		Key:         *req.Key,
		IsActive:    *req.IsActive,
		Code:        utils.RandStringBytes(10, true),
		Description: *utils.SafeStringPointer(req.Description, ""),
	}
}
