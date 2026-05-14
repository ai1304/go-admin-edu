package models

type EduExpert struct {
	BaseModel
	Name         string `json:"name" gorm:"size:128;not null;index;comment:专家姓名"`
	Title        string `json:"title" gorm:"size:128;comment:职称"`
	Organization string `json:"organization" gorm:"size:255;comment:机构"`
	AvatarFileId int    `json:"avatarFileId" gorm:"index;comment:头像文件ID"`
	Specialties  string `json:"specialties" gorm:"size:512;comment:擅长领域"`
	Introduction string `json:"introduction" gorm:"type:text;comment:简介"`
	Status       int    `json:"status" gorm:"default:1;comment:状态"`
}

func (*EduExpert) TableName() string {
	return "edu_expert"
}

type EduExpertResource struct {
	BaseModel
	ExpertId   int    `json:"expertId" gorm:"index;comment:专家ID"`
	Title      string `json:"title" gorm:"size:255;not null;comment:资源标题"`
	Type       string `json:"type" gorm:"size:32;index;comment:资源类型"`
	ResourceId int    `json:"resourceId" gorm:"index;comment:关联资源ID"`
	CourseId   int    `json:"courseId" gorm:"index;comment:关联课程ID"`
	FileId     int    `json:"fileId" gorm:"index;comment:文件ID"`
	Status     int    `json:"status" gorm:"default:1;comment:状态"`
}

func (*EduExpertResource) TableName() string {
	return "edu_expert_resource"
}
