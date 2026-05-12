package models

type EduRegion struct {
	BaseModel
	Name     string `json:"name" gorm:"size:128;not null;comment:区域名称"`
	Code     string `json:"code" gorm:"size:64;index;comment:区域编码"`
	ParentId int    `json:"parentId" gorm:"index;comment:上级区域ID"`
	Sort     int    `json:"sort" gorm:"comment:排序"`
	Status   int    `json:"status" gorm:"default:1;comment:状态"`
	Remark   string `json:"remark" gorm:"size:512;comment:备注"`
}

func (*EduRegion) TableName() string {
	return "edu_region"
}
