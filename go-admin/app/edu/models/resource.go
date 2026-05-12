package models

const (
	ResourceStatusDraft     = "draft"
	ResourceStatusReviewing = "reviewing"
	ResourceStatusPublished = "published"
	ResourceStatusRejected  = "rejected"
	ResourceStatusOffline   = "offline"
)

type EduResource struct {
	BaseModel
	Title            string `json:"title" gorm:"size:255;not null;index;comment:资源标题"`
	Summary          string `json:"summary" gorm:"size:1024;comment:资源简介"`
	CoverFileId      int    `json:"coverFileId" gorm:"index;comment:封面文件ID"`
	StageCategoryId  int    `json:"stageCategoryId" gorm:"index;comment:学段分类ID"`
	DisabilityTypeId int    `json:"disabilityTypeId" gorm:"index;comment:障碍类型分类ID"`
	ResourceTypeId   int    `json:"resourceTypeId" gorm:"index;comment:资源类型分类ID"`
	AbilityDomainId  int    `json:"abilityDomainId" gorm:"index;comment:能力领域分类ID"`
	TopicCategoryId  int    `json:"topicCategoryId" gorm:"index;comment:专题分类ID"`
	AuthorName       string `json:"authorName" gorm:"size:128;comment:作者/教师"`
	Keywords         string `json:"keywords" gorm:"size:512;comment:关键词"`
	Status           string `json:"status" gorm:"size:32;default:draft;index;comment:状态"`
	ViewCount        int64  `json:"viewCount" gorm:"default:0;comment:浏览量"`
	DownloadCount    int64  `json:"downloadCount" gorm:"default:0;comment:下载量"`
	FavoriteCount    int64  `json:"favoriteCount" gorm:"default:0;comment:收藏量"`
}

func (*EduResource) TableName() string {
	return "edu_resource"
}
