package model

import (
	"github.com/google/uuid"
	"good-wedding/pkg/model/paging"
)

type Todo struct {
	BaseModel
	Name        string `json:"name" gorm:"column:name;type:text;not null"`
	Key         string `json:"key" gorm:"column:key;type:text;unique;not null"`
	IsActive    bool   `json:"is_active" gorm:"column:is_active;not null"`
	Code        string `json:"code" gorm:"column:code;type:text;unique;not null"`
	Description string `json:"description" gorm:"column:description;type:text;default:null"`
}

func (m *Todo) TableName() string {
	return "todo"
}

type CreateTodoRequest struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Name        *string    `json:"name" valid:"Required"`
	Key         *string    `json:"key" valid:"Required"`
	IsActive    *bool      `json:"is_active" valid:"Required"`
	Code        *string    `json:"code" valid:"Required"`
	Description *string    `json:"description"`
}

type CreateTodoResponse struct {
	Meta *MetaData              `json:"meta"`
	Data CreateTodoDataResponse `json:"data"`
}

type CreateTodoDataResponse struct {
	CreatorID uuid.UUID `json:"creator_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	IsActive  bool      `json:"is_active"`
	Code      string    `json:"code"`
}

type TodoFilterRequest struct {
	FromDate  *int64  `json:"from_date" form:"from_date"`
	ToDate    *int64  `json:"to_date" form:"to_date"`
	CreatorID *string `json:"creator_id" form:"creator_id"`
	Name      *string `json:"name" form:"name"`
	Key       *string `json:"key" form:"key"`
	IsActive  *bool   `json:"is_active" form:"is_active"`
	Code      *string `json:"code" form:"code"`
}

type TodoFilter struct {
	TodoFilterRequest
	Pager *paging.Pager
}

type TodoFilterResult struct {
	Filter  *TodoFilter `json:"filter"`
	Records []*Todo     `json:"data"`
}

type TodoUpdateRequest struct {
	Name        *string `json:"name"`
	Key         *string `json:"key"`
	IsActive    *bool   `json:"is_active"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
}

type TodoUpdateResponse struct {
	Meta *MetaData              `json:"meta"`
	Data TodoUpdateDataResponse `json:"data"`
}

type TodoUpdateDataResponse struct {
	Name        *string `json:"name"`
	Key         *string `json:"key"`
	IsActive    *bool   `json:"is_active"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
}

type TodoKafkaMessage struct {
	Payload struct {
		Before *Todo `json:"before"`
		After  *Todo `json:"after"`
	} `json:"payload"`
}
