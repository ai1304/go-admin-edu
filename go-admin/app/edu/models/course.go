package models

type EduCourse struct {
	BaseModel
	Title            string `json:"title" gorm:"size:255;not null;index;comment:课程标题"`
	Summary          string `json:"summary" gorm:"size:1024;comment:课程简介"`
	CoverFileId      int    `json:"coverFileId" gorm:"index;comment:封面文件ID"`
	VideoFileId      int    `json:"videoFileId" gorm:"index;comment:课程视频文件ID"`
	StageCategoryId  int    `json:"stageCategoryId" gorm:"index;comment:学段分类ID"`
	DisabilityTypeId int    `json:"disabilityTypeId" gorm:"index;comment:障碍类型分类ID"`
	Category         string `json:"category" gorm:"size:128;index;comment:课程分类"`
	Difficulty       string `json:"difficulty" gorm:"size:32;comment:难度"`
	TeacherName      string `json:"teacherName" gorm:"size:128;comment:教师名称"`
	Organization     string `json:"organization" gorm:"size:255;comment:机构"`
	Objectives       string `json:"objectives" gorm:"type:text;comment:教学目标"`
	Status           string `json:"status" gorm:"size:32;default:draft;index;comment:状态"`
	ViewCount        int64  `json:"viewCount" gorm:"default:0;comment:浏览量"`
	LearnerCount     int64  `json:"learnerCount" gorm:"default:0;comment:学习人数"`
	Sort             int    `json:"sort" gorm:"default:0;comment:排序"`
}

func (*EduCourse) TableName() string {
	return "edu_course"
}

type EduCourseChapter struct {
	BaseModel
	CourseId int    `json:"courseId" gorm:"index;comment:课程ID"`
	Title    string `json:"title" gorm:"size:255;not null;comment:章节标题"`
	Sort     int    `json:"sort" gorm:"comment:排序"`
	Status   int    `json:"status" gorm:"default:1;comment:状态"`
}

func (*EduCourseChapter) TableName() string {
	return "edu_course_chapter"
}

type EduCourseLesson struct {
	BaseModel
	CourseId        int    `json:"courseId" gorm:"index;comment:课程ID"`
	ChapterId       int    `json:"chapterId" gorm:"index;comment:章节ID"`
	Title           string `json:"title" gorm:"size:255;not null;comment:课时标题"`
	VideoFileId     int    `json:"videoFileId" gorm:"index;comment:视频文件ID"`
	DurationSeconds int    `json:"durationSeconds" gorm:"comment:时长秒"`
	Sort            int    `json:"sort" gorm:"comment:排序"`
	Status          int    `json:"status" gorm:"default:1;comment:状态"`
}

func (*EduCourseLesson) TableName() string {
	return "edu_course_lesson"
}

type EduLearningRecord struct {
	BaseModel
	CourseId       int    `json:"courseId" gorm:"index;comment:课程ID"`
	LessonId       int    `json:"lessonId" gorm:"index;comment:课时ID"`
	UserId         int    `json:"userId" gorm:"index;comment:学习用户ID"`
	ClientKey      string `json:"clientKey" gorm:"size:128;index;comment:客户端标识"`
	Progress       int    `json:"progress" gorm:"default:0;comment:进度百分比"`
	WatchedSeconds int    `json:"watchedSeconds" gorm:"default:0;comment:观看秒数"`
	Status         string `json:"status" gorm:"size:32;default:learning;comment:学习状态"`
}

func (*EduLearningRecord) TableName() string {
	return "edu_learning_record"
}

type EduAssignment struct {
	BaseModel
	CourseId int    `json:"courseId" gorm:"index;comment:课程ID"`
	Title    string `json:"title" gorm:"size:255;not null;comment:作业标题"`
	Content  string `json:"content" gorm:"type:text;comment:作业内容"`
	Status   int    `json:"status" gorm:"default:1;comment:状态"`
}

func (*EduAssignment) TableName() string {
	return "edu_assignment"
}

type EduAssignmentSubmission struct {
	BaseModel
	AssignmentId int    `json:"assignmentId" gorm:"index;comment:作业ID"`
	CourseId     int    `json:"courseId" gorm:"index;comment:课程ID"`
	UserId       int    `json:"userId" gorm:"index;comment:提交用户ID"`
	ClientKey    string `json:"clientKey" gorm:"size:128;index;comment:客户端标识"`
	Nickname     string `json:"nickname" gorm:"size:128;comment:昵称"`
	Content      string `json:"content" gorm:"type:text;comment:提交内容"`
	FileId       int    `json:"fileId" gorm:"index;comment:附件文件ID"`
	Score        int    `json:"score" gorm:"comment:分数"`
	Status       string `json:"status" gorm:"size:32;default:submitted;comment:状态"`
}

func (*EduAssignmentSubmission) TableName() string {
	return "edu_assignment_submission"
}
