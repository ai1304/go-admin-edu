package models

type EduResourceFile struct {
	BaseModel
	ResourceId   int    `json:"resourceId" gorm:"index;comment:资源ID"`
	OriginalName string `json:"originalName" gorm:"size:255;comment:原始文件名"`
	ObjectKey    string `json:"objectKey" gorm:"size:512;not null;comment:对象存储Key"`
	Bucket       string `json:"bucket" gorm:"size:128;comment:存储桶"`
	ContentType  string `json:"contentType" gorm:"size:128;comment:文件类型"`
	Ext          string `json:"ext" gorm:"size:32;comment:扩展名"`
	Size         int64  `json:"size" gorm:"comment:文件大小"`
	Hash         string `json:"hash" gorm:"size:128;index;comment:文件Hash"`
	Usage        string `json:"usage" gorm:"size:32;index;comment:用途"`
}

func (*EduResourceFile) TableName() string {
	return "edu_resource_file"
}
