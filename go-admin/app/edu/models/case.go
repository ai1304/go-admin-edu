package models

type EduCase struct {
	BaseModel
	Title          string `json:"title" gorm:"size:255;not null;index;comment:案例名称"`
	StudentName    string `json:"studentName" gorm:"size:128;comment:学生姓名"`
	StudentCode    string `json:"studentCode" gorm:"size:64;index;comment:学生编号"`
	Gender         string `json:"gender" gorm:"size:16;comment:性别"`
	Birthday       string `json:"birthday" gorm:"size:32;comment:生日"`
	DisabilityType string `json:"disabilityType" gorm:"size:128;comment:障碍类型"`
	Summary        string `json:"summary" gorm:"size:1024;comment:案例摘要"`
	Status         string `json:"status" gorm:"size:32;default:draft;index;comment:状态"`
}

func (*EduCase) TableName() string {
	return "edu_case"
}

type EduCaseIEP struct {
	BaseModel
	CaseId     int    `json:"caseId" gorm:"index;comment:案例ID"`
	Title      string `json:"title" gorm:"size:255;not null;comment:IEP标题"`
	Goal       string `json:"goal" gorm:"type:text;comment:目标"`
	Plan       string `json:"plan" gorm:"type:text;comment:计划"`
	Evaluation string `json:"evaluation" gorm:"type:text;comment:评价"`
	Status     string `json:"status" gorm:"size:32;default:draft;comment:状态"`
}

func (*EduCaseIEP) TableName() string {
	return "edu_case_iep"
}

type EduCaseAssessment struct {
	BaseModel
	CaseId     int    `json:"caseId" gorm:"index;comment:案例ID"`
	ToolName   string `json:"toolName" gorm:"size:255;comment:评估工具"`
	Result     string `json:"result" gorm:"type:text;comment:评估结果"`
	AssessedAt string `json:"assessedAt" gorm:"size:32;comment:评估时间"`
}

func (*EduCaseAssessment) TableName() string {
	return "edu_case_assessment"
}

type EduCaseIntervention struct {
	BaseModel
	CaseId    int    `json:"caseId" gorm:"index;comment:案例ID"`
	Title     string `json:"title" gorm:"size:255;not null;comment:干预方案标题"`
	Content   string `json:"content" gorm:"type:text;comment:干预内容"`
	StartDate string `json:"startDate" gorm:"size:32;comment:开始日期"`
	EndDate   string `json:"endDate" gorm:"size:32;comment:结束日期"`
	Status    string `json:"status" gorm:"size:32;default:active;comment:状态"`
}

func (*EduCaseIntervention) TableName() string {
	return "edu_case_intervention"
}
