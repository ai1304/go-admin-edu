package models

type EduResourceReview struct {
	BaseModel
	ResourceId   int    `json:"resourceId" gorm:"index;comment:资源ID"`
	Action       string `json:"action" gorm:"size:32;index;comment:审核动作"`
	Comment      string `json:"comment" gorm:"size:1024;comment:审核意见"`
	BeforeStatus string `json:"beforeStatus" gorm:"size:32;comment:原状态"`
	AfterStatus  string `json:"afterStatus" gorm:"size:32;comment:新状态"`
}

func (*EduResourceReview) TableName() string {
	return "edu_resource_review"
}
