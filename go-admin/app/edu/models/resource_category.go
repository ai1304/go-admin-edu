package models

type EduResourceCategory struct {
	BaseModel
	Name     string `json:"name" gorm:"size:128;not null;comment:分类名称"`
	Code     string `json:"code" gorm:"size:64;index;comment:分类编码"`
	Type     string `json:"type" gorm:"size:32;index;comment:分类类型"`
	ParentId int    `json:"parentId" gorm:"index;comment:上级分类ID"`
	Sort     int    `json:"sort" gorm:"comment:排序"`
	Status   int    `json:"status" gorm:"default:1;comment:状态"`
	Remark   string `json:"remark" gorm:"size:512;comment:备注"`
}

func (*EduResourceCategory) TableName() string {
	return "edu_resource_category"
}
