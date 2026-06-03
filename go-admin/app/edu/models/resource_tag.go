package models

type EduResourceTag struct {
	BaseModel
	Name   string `json:"name" gorm:"size:64;not null;index;comment:标签名称"`
	Status int    `json:"status" gorm:"default:1;comment:状态"`
}

func (*EduResourceTag) TableName() string {
	return "edu_resource_tag"
}

type EduResourceTagRelation struct {
	BaseModel
	ResourceId int `json:"resourceId" gorm:"index;comment:资源ID"`
	TagId      int `json:"tagId" gorm:"index;comment:标签ID"`
}

func (*EduResourceTagRelation) TableName() string {
	return "edu_resource_tag_relation"
}
