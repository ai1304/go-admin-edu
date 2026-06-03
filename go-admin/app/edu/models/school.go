package models

type EduSchool struct {
	BaseModel
	Name    string `json:"name" gorm:"size:128;not null;comment:学校名称"`
	Code    string `json:"code" gorm:"size:64;index;comment:学校编码"`
	Address string `json:"address" gorm:"size:255;comment:学校地址"`
	Contact string `json:"contact" gorm:"size:64;comment:联系人"`
	Phone   string `json:"phone" gorm:"size:32;comment:联系电话"`
	Status  int    `json:"status" gorm:"default:1;comment:状态"`
	Remark  string `json:"remark" gorm:"size:512;comment:备注"`
}

func (*EduSchool) TableName() string {
	return "edu_school"
}
