package models

const (
	NewsStatusDraft     = "draft"
	NewsStatusPublished = "published"
	NewsStatusOffline   = "offline"
)

type EduNews struct {
	BaseModel
	Title       string `json:"title" gorm:"size:255;not null;index;comment:资讯标题"`
	ModuleType  string `json:"moduleType" gorm:"size:32;index;comment:资讯类型"`
	Source      string `json:"source" gorm:"size:128;comment:来源"`
	CoverFileId int    `json:"coverFileId" gorm:"index;comment:封面文件ID"`
	CoverUrl    string `json:"coverUrl" gorm:"size:512;comment:外部封面地址"`
	Summary     string `json:"summary" gorm:"size:1024;comment:摘要"`
	Content     string `json:"content" gorm:"type:longtext;comment:正文"`
	Keywords    string `json:"keywords" gorm:"size:255;comment:关键词"`
	PublishTime string `json:"publishTime" gorm:"size:32;index;comment:发布时间"`
	ViewCount   int64  `json:"viewCount" gorm:"default:0;comment:浏览量"`
	LikeCount   int64  `json:"likeCount" gorm:"default:0;comment:点赞量"`
	Status      string `json:"status" gorm:"size:32;default:draft;index;comment:状态"`
	IsTop       int    `json:"isTop" gorm:"default:0;comment:是否置顶"`
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`
}

func (*EduNews) TableName() string {
	return "edu_news"
}
