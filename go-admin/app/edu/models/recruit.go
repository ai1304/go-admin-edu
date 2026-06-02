package models

type EduRecruitCompany struct {
	BaseModel
	CompanyName       string `json:"companyName" gorm:"size:255;not null;index;comment:企业名称"`
	CreditCode        string `json:"creditCode" gorm:"size:64;not null;index;comment:统一社会信用代码"`
	CompanyNature     string `json:"companyNature" gorm:"size:64;index;comment:企业性质"`
	Industry          string `json:"industry" gorm:"size:128;index;comment:所属行业"`
	CompanySize       string `json:"companySize" gorm:"size:64;index;comment:企业规模"`
	Region            string `json:"region" gorm:"size:128;index;comment:所在地区"`
	Address           string `json:"address" gorm:"size:255;comment:详细地址"`
	Website           string `json:"website" gorm:"size:512;comment:官网链接"`
	LogoUrl           string `json:"logoUrl" gorm:"size:512;comment:企业Logo"`
	ContactName       string `json:"contactName" gorm:"size:64;index;comment:联系人"`
	ContactTitle      string `json:"contactTitle" gorm:"size:64;comment:联系人职务"`
	ContactPhone      string `json:"contactPhone" gorm:"size:32;comment:联系电话"`
	ContactEmail      string `json:"contactEmail" gorm:"size:128;comment:联系邮箱"`
	Intro             string `json:"intro" gorm:"type:text;comment:企业简介"`
	MainBusiness      string `json:"mainBusiness" gorm:"type:text;comment:主营业务"`
	TalentNeeds       string `json:"talentNeeds" gorm:"size:512;comment:人才需求方向"`
	Cooperation       string `json:"cooperation" gorm:"type:text;comment:校企合作方向"`
	LicenseMaterial   string `json:"licenseMaterial" gorm:"type:text;comment:营业执照或单位证明"`
	Qualification     string `json:"qualification" gorm:"type:text;comment:企业资质材料"`
	AuthorizationNote string `json:"authorizationNote" gorm:"type:text;comment:联系人授权说明"`
	Tags              string `json:"tags" gorm:"size:512;comment:企业标签"`
	CertStatus        string `json:"certStatus" gorm:"size:32;default:pending;index;comment:认证状态"`
	Status            string `json:"status" gorm:"size:32;default:pending;index;comment:入驻状态"`
	ReviewStatus      string `json:"reviewStatus" gorm:"size:32;default:pending;index;comment:审核状态"`
	ReviewOpinion     string `json:"reviewOpinion" gorm:"type:text;comment:审核意见"`
	ReviewedAt        string `json:"reviewedAt" gorm:"size:32;comment:审核时间"`
	JobCount          int64  `json:"jobCount" gorm:"-"`
}

func (*EduRecruitCompany) TableName() string {
	return "edu_recruit_company"
}

type EduRecruitJob struct {
	BaseModel
	CompanyId        int    `json:"companyId" gorm:"index;comment:所属企业ID"`
	CompanyName      string `json:"companyName" gorm:"size:255;index;comment:企业名称"`
	Industry         string `json:"industry" gorm:"size:128;index;comment:所属行业"`
	JobName          string `json:"jobName" gorm:"size:255;not null;index;comment:岗位名称"`
	JobType          string `json:"jobType" gorm:"size:64;index;comment:岗位类型"`
	MajorDirection   string `json:"majorDirection" gorm:"size:128;index;comment:专业方向"`
	Location         string `json:"location" gorm:"size:128;index;comment:工作地点"`
	Headcount        int    `json:"headcount" gorm:"default:1;comment:招聘人数"`
	SalaryRange      string `json:"salaryRange" gorm:"size:128;index;comment:薪资范围"`
	Education        string `json:"education" gorm:"size:64;index;comment:学历要求"`
	MajorRequirement string `json:"majorRequirement" gorm:"size:255;comment:专业要求"`
	Experience       string `json:"experience" gorm:"size:128;comment:工作经验要求"`
	RecruitTarget    string `json:"recruitTarget" gorm:"size:128;comment:招聘对象"`
	PublishTime      string `json:"publishTime" gorm:"size:32;index;comment:发布时间"`
	Deadline         string `json:"deadline" gorm:"size:32;index;comment:截止时间"`
	Responsibilities string `json:"responsibilities" gorm:"type:text;comment:岗位职责"`
	Requirements     string `json:"requirements" gorm:"type:text;comment:任职要求"`
	WorkTime         string `json:"workTime" gorm:"type:text;comment:工作时间"`
	Benefits         string `json:"benefits" gorm:"type:text;comment:福利待遇"`
	Training         string `json:"training" gorm:"type:text;comment:培养机制"`
	OtherNotes       string `json:"otherNotes" gorm:"type:text;comment:其他说明"`
	Tags             string `json:"tags" gorm:"size:512;comment:岗位标签"`
	ContactName      string `json:"contactName" gorm:"size:64;comment:岗位联系人"`
	ContactTitle     string `json:"contactTitle" gorm:"size:64;comment:联系人职务"`
	ContactPhone     string `json:"contactPhone" gorm:"size:32;comment:联系电话"`
	ContactEmail     string `json:"contactEmail" gorm:"size:128;comment:联系邮箱"`
	ContactAddress   string `json:"contactAddress" gorm:"size:255;comment:联系地址"`
	ExternalLink     string `json:"externalLink" gorm:"size:512;comment:岗位外部链接"`
	JobAttachment    string `json:"jobAttachment" gorm:"type:text;comment:岗位说明附件"`
	RecruitBrochure  string `json:"recruitBrochure" gorm:"type:text;comment:企业招聘简章"`
	Status           string `json:"status" gorm:"size:32;default:pending;index;comment:岗位状态"`
	ReviewStatus     string `json:"reviewStatus" gorm:"size:32;default:pending;index;comment:审核状态"`
	ReviewOpinion    string `json:"reviewOpinion" gorm:"type:text;comment:审核意见"`
	ReviewedAt       string `json:"reviewedAt" gorm:"size:32;comment:审核时间"`
}

func (*EduRecruitJob) TableName() string {
	return "edu_recruit_job"
}

type EduRecruitReview struct {
	BaseModel
	TargetType string `json:"targetType" gorm:"size:32;index;comment:审核对象类型"`
	TargetId   int    `json:"targetId" gorm:"index;comment:审核对象ID"`
	Action     string `json:"action" gorm:"size:32;index;comment:审核动作"`
	Status     string `json:"status" gorm:"size:32;index;comment:审核状态"`
	Opinion    string `json:"opinion" gorm:"type:text;comment:审核意见"`
	Reason     string `json:"reason" gorm:"size:128;comment:驳回原因"`
}

func (*EduRecruitReview) TableName() string {
	return "edu_recruit_review"
}
