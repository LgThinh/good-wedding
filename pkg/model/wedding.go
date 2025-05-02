package model

import "github.com/google/uuid"

type ObjectMedia struct {
	BaseModel
	Name       string `json:"name" gorm:"column:name;type:text;NOT NULL"`
	ObjectType string `json:"object_type" gorm:"column:object_type;type:text;NOT NULL"`
	FileType   string `json:"file_type" gorm:"column:file_type;type:text;NOT NULL"`
	Url        string `json:"url" gorm:"column:url;type:text;NOT NULL"`
	Size       int64  `json:"size" gorm:"column:size;NOT NULL"`
}

func (m ObjectMedia) TableName() string {
	return "object_media"
}

type User struct {
	BaseModel
	UserName string ` json:"username" gorm:"column:username"`
}

func (m User) TableName() string {
	return "user"
}

type Comment struct {
	BaseModel
	ObjectID uuid.UUID `json:"object_id" gorm:"column:object_id;type:uuid;NOT NULL"`
	UserID   uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;NOT NULL"`
	Comment  string    `json:"comment" gorm:"column:comment;type:text;NOT NULL"`
}

func (m Comment) TableName() string {
	return "comment"
}

type WeddingWish struct {
	BaseModel
	UserID  uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;NOT NULL"`
	Comment string    `json:"comment" gorm:"column:comment;type:text;NOT NULL"`
}

func (m WeddingWish) TableName() string {
	return "wedding_wish"
}
