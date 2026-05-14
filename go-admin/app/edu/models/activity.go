package models

type EduActivity struct {
	BaseModel
	Title       string `json:"title" gorm:"size:255;not null;index;comment:活动名称"`
	Summary     string `json:"summary" gorm:"size:1024;comment:活动简介"`
	StartTime   string `json:"startTime" gorm:"size:32;index;comment:开始时间"`
	EndTime     string `json:"endTime" gorm:"size:32;comment:结束时间"`
	Location    string `json:"location" gorm:"size:255;comment:地点"`
	Organizer   string `json:"organizer" gorm:"size:128;comment:主办方"`
	Status      string `json:"status" gorm:"size:32;default:draft;index;comment:状态"`
	SignupCount int64  `json:"signupCount" gorm:"default:0;comment:报名人数"`
}

func (*EduActivity) TableName() string {
	return "edu_activity"
}

type EduActivitySignup struct {
	BaseModel
	ActivityId int    `json:"activityId" gorm:"index;comment:活动ID"`
	UserId     int    `json:"userId" gorm:"index;comment:用户ID"`
	Name       string `json:"name" gorm:"size:128;comment:报名姓名"`
	Phone      string `json:"phone" gorm:"size:32;comment:电话"`
	Status     string `json:"status" gorm:"size:32;default:signed;comment:状态"`
}

func (*EduActivitySignup) TableName() string {
	return "edu_activity_signup"
}

type EduActivityCheckin struct {
	BaseModel
	ActivityId int    `json:"activityId" gorm:"index;comment:活动ID"`
	UserId     int    `json:"userId" gorm:"index;comment:用户ID"`
	CheckinAt  string `json:"checkinAt" gorm:"size:32;comment:签到时间"`
	Status     string `json:"status" gorm:"size:32;default:checked;comment:状态"`
}

func (*EduActivityCheckin) TableName() string {
	return "edu_activity_checkin"
}

type EduActivityOutcome struct {
	BaseModel
	ActivityId int    `json:"activityId" gorm:"index;comment:活动ID"`
	Title      string `json:"title" gorm:"size:255;not null;comment:成果标题"`
	Content    string `json:"content" gorm:"type:text;comment:成果内容"`
	FileId     int    `json:"fileId" gorm:"index;comment:附件文件ID"`
	Status     int    `json:"status" gorm:"default:1;comment:状态"`
}

func (*EduActivityOutcome) TableName() string {
	return "edu_activity_outcome"
}
