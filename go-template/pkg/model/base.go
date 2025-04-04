package model

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Pagination struct {
	Page     int
	PageSize int
}

type UriParse struct {
	ID []string `json:"id" uri:"id"`
}

// MsMetadata describes the structure.
type BaseModel struct {
	ID        uuid.UUID       `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatorID *uuid.UUID      `json:"creator_id,omitempty"`
	UpdaterID *uuid.UUID      `json:"updater_id,omitempty"`
	CreatedAt time.Time       `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type MetaData struct {
	TraceID string `json:"traceId"`
	Success bool   `json:"success"`
}

func NewMetaData() *MetaData {
	return &MetaData{
		TraceID: "your-trace-id",
		Success: true,
	}
}

func NewMetaDataWithTraceID(ctx context.Context) *MetaData {
	requestID, ok := ctx.Value("x-request-id").(string)
	if !ok {
		requestID = "unknown-trace-id"
	}
	return &MetaData{
		TraceID: requestID,
		Success: true,
	}
}
