package model

type Object struct {
	BaseModel
	Name       string `json:"name" gorm:"column:name;type:text;NOT NULL"`
	ObjectType string `json:"object_type" gorm:"column:object_type;type:text;NOT NULL"`
	Url        string `json:"url" gorm:"column:url;type:text;NOT NULL"`
	Size       int64  `json:"size" gorm:"column:size;NOT NULL"`
}

func (m Object) TableName() string {
	return "object"
}
