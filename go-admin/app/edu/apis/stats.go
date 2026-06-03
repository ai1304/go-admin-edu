package apis

import (
	"encoding/csv"
	"fmt"
	"go-admin/app/edu/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
)

type EduStats struct {
	api.Api
}

type statsPair struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type statsTrend struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

type schoolContribution struct {
	SchoolId      int    `json:"schoolId"`
	SchoolName    string `json:"schoolName"`
	ResourceCount int64  `json:"resourceCount"`
	CourseCount   int64  `json:"courseCount"`
	ActivityCount int64  `json:"activityCount"`
	Total         int64  `json:"total"`
}

type teacherActivity struct {
	UserId        int   `json:"userId"`
	ResourceCount int64 `json:"resourceCount"`
	CourseCount   int64 `json:"courseCount"`
	CaseCount     int64 `json:"caseCount"`
	Total         int64 `json:"total"`
}

type studentLearning struct {
	Identity        string `json:"identity"`
	UserId          int    `json:"userId"`
	ClientKey       string `json:"clientKey"`
	LessonCount     int64  `json:"lessonCount"`
	FinishedCount   int64  `json:"finishedCount"`
	SubmissionCount int64  `json:"submissionCount"`
}

func (e EduStats) Overview(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	counts := gin.H{}
	counts["regions"] = e.countModel(&models.EduRegion{})
	counts["schools"] = e.countModel(&models.EduSchool{})
	counts["resources"] = e.countModel(&models.EduResource{})
	counts["publishedResources"] = e.countWhere(&models.EduResource{}, "status = ?", models.ResourceStatusPublished)
	counts["courses"] = e.countModel(&models.EduCourse{})
	counts["publishedCourses"] = e.countWhere(&models.EduCourse{}, "status = ?", "published")
	counts["activities"] = e.countModel(&models.EduActivity{})
	counts["publishedActivities"] = e.countWhere(&models.EduActivity{}, "status = ?", "published")
	counts["cases"] = e.countModel(&models.EduCase{})
	counts["experts"] = e.countModel(&models.EduExpert{})
	counts["downloads"] = e.sumInt64(&models.EduResource{}, "download_count")
	counts["views"] = e.sumInt64(&models.EduResource{}, "view_count")
	counts["learners"] = e.countLearningIdentities()
	e.OK(counts, "查询成功")
}

func (e EduStats) Resources(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(gin.H{
		"byStatus":       e.groupCount("edu_resource", "status", "status <> ''"),
		"byResourceType": e.groupCategory("resource_type_id"),
		"uploadTrend":    e.trend("edu_resource", 14),
		"topViewed":      e.topResources("view_count"),
		"topDownloaded":  e.topResources("download_count"),
		"topFavorited":   e.topResources("favorite_count"),
	}, "查询成功")
}

func (e EduStats) Courses(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(gin.H{
		"byStatus":        e.groupCount("edu_course", "status", "status <> ''"),
		"learningTrend":   e.trend("edu_learning_record", 14),
		"totalLearners":   e.countLearningIdentities(),
		"finishedLessons": e.countWhere(&models.EduLearningRecord{}, "status = ?", "finished"),
		"submissions":     e.countModel(&models.EduAssignmentSubmission{}),
		"topCourses":      e.topCourses(),
	}, "查询成功")
}

func (e EduStats) Activities(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(gin.H{
		"byStatus":    e.groupCount("edu_activity", "status", "status <> ''"),
		"signups":     e.countModel(&models.EduActivitySignup{}),
		"checkins":    e.countModel(&models.EduActivityCheckin{}),
		"outcomes":    e.countModel(&models.EduActivityOutcome{}),
		"topActivity": e.topActivities(),
	}, "查询成功")
}

func (e EduStats) Schools(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]schoolContribution, 0)
	_ = e.Orm.Raw(`
		select s.id as school_id, s.name as school_name,
			count(distinct r.id) as resource_count,
			count(distinct c.id) as course_count,
			count(distinct a.id) as activity_count,
			count(distinct r.id) + count(distinct c.id) + count(distinct a.id) as total
		from edu_school s
		left join edu_resource r on r.school_id = s.id and r.deleted_at is null
		left join edu_course c on c.school_id = s.id and c.deleted_at is null
		left join edu_activity a on a.school_id = s.id and a.deleted_at is null
		where s.deleted_at is null
		group by s.id, s.name
		order by total desc, s.id desc
		limit 20
	`).Scan(&list).Error
	e.OK(list, "查询成功")
}

func (e EduStats) Teachers(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]teacherActivity, 0)
	_ = e.Orm.Raw(`
		select u.user_id,
			sum(u.resource_count) as resource_count,
			sum(u.course_count) as course_count,
			sum(u.case_count) as case_count,
			sum(u.resource_count + u.course_count + u.case_count) as total
		from (
			select create_by as user_id, count(*) as resource_count, 0 as course_count, 0 as case_count from edu_resource where deleted_at is null and create_by <> 0 group by create_by
			union all
			select create_by as user_id, 0 as resource_count, count(*) as course_count, 0 as case_count from edu_course where deleted_at is null and create_by <> 0 group by create_by
			union all
			select create_by as user_id, 0 as resource_count, 0 as course_count, count(*) as case_count from edu_case where deleted_at is null and create_by <> 0 group by create_by
		) u
		group by u.user_id
		order by total desc
		limit 20
	`).Scan(&list).Error
	e.OK(list, "查询成功")
}

func (e EduStats) Students(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]studentLearning, 0)
	_ = e.Orm.Raw(`
		select l.identity, l.user_id, l.client_key, l.lesson_count, l.finished_count, coalesce(s.submission_count, 0) as submission_count
		from (
			select
				case when user_id <> 0 then concat('u:', user_id) else concat('c:', client_key) end as identity,
				max(user_id) as user_id,
				max(client_key) as client_key,
				count(*) as lesson_count,
				sum(case when status = 'finished' then 1 else 0 end) as finished_count
			from edu_learning_record
			where deleted_at is null and (user_id <> 0 or client_key <> '')
			group by identity
		) l
		left join (
			select
				case when user_id <> 0 then concat('u:', user_id) else concat('c:', client_key) end as identity,
				count(*) as submission_count
			from edu_assignment_submission
			where deleted_at is null and (user_id <> 0 or client_key <> '')
			group by identity
		) s on s.identity = l.identity
		order by l.lesson_count desc
		limit 20
	`).Scan(&list).Error
	e.OK(list, "查询成功")
}

func (e EduStats) Cases(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(gin.H{
		"byStatus":       e.groupCount("edu_case", "status", "status <> ''"),
		"ieps":           e.countModel(&models.EduCaseIEP{}),
		"assessments":    e.countModel(&models.EduCaseAssessment{}),
		"interventions":  e.countModel(&models.EduCaseIntervention{}),
		"accessLogs":     e.countModel(&models.EduCaseAccessLog{}),
		"authorizations": e.countModel(&models.EduCaseAuthorization{}),
	}, "查询成功")
}

func (e EduStats) Export(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	filename := fmt.Sprintf("edu-stats-%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()
	_ = writer.Write([]string{"Metric", "Value"})
	rows := []statsPair{
		{"regions", e.countModel(&models.EduRegion{})},
		{"schools", e.countModel(&models.EduSchool{})},
		{"resources", e.countModel(&models.EduResource{})},
		{"published_resources", e.countWhere(&models.EduResource{}, "status = ?", models.ResourceStatusPublished)},
		{"courses", e.countModel(&models.EduCourse{})},
		{"published_courses", e.countWhere(&models.EduCourse{}, "status = ?", "published")},
		{"activities", e.countModel(&models.EduActivity{})},
		{"cases", e.countModel(&models.EduCase{})},
		{"experts", e.countModel(&models.EduExpert{})},
		{"learners", e.countLearningIdentities()},
	}
	for _, row := range rows {
		_ = writer.Write([]string{row.Name, strconv.FormatInt(row.Value, 10)})
	}
}

func (e EduStats) countModel(model interface{}) int64 {
	var count int64
	_ = e.Orm.Model(model).Count(&count).Error
	return count
}

func (e EduStats) countWhere(model interface{}, query string, args ...interface{}) int64 {
	var count int64
	_ = e.Orm.Model(model).Where(query, args...).Count(&count).Error
	return count
}

func (e EduStats) sumInt64(model interface{}, column string) int64 {
	var value int64
	_ = e.Orm.Model(model).Select("coalesce(sum(" + column + "), 0)").Scan(&value).Error
	return value
}

func (e EduStats) countLearningIdentities() int64 {
	var count int64
	_ = e.Orm.Raw(`
		select count(*) from (
			select case when user_id <> 0 then concat('u:', user_id) else concat('c:', client_key) end as identity
			from edu_learning_record
			where deleted_at is null and (user_id <> 0 or client_key <> '')
			group by identity
		) t
	`).Scan(&count).Error
	return count
}

func (e EduStats) groupCount(table string, column string, condition string) []statsPair {
	list := make([]statsPair, 0)
	sql := fmt.Sprintf("select %s as name, count(*) as value from %s where deleted_at is null and %s group by %s order by value desc", column, table, condition, column)
	_ = e.Orm.Raw(sql).Scan(&list).Error
	return list
}

func (e EduStats) groupCategory(column string) []statsPair {
	list := make([]statsPair, 0)
	sql := fmt.Sprintf(`
		select coalesce(c.name, '未分类') as name, count(r.id) as value
		from edu_resource r
		left join edu_resource_category c on c.id = r.%s and c.deleted_at is null
		where r.deleted_at is null
		group by name
		order by value desc
	`, column)
	_ = e.Orm.Raw(sql).Scan(&list).Error
	return list
}

func (e EduStats) trend(table string, days int) []statsTrend {
	list := make([]statsTrend, 0)
	if days <= 0 {
		days = 14
	}
	sql := fmt.Sprintf(`
		select date(created_at) as date, count(*) as value
		from %s
		where deleted_at is null and created_at >= date_sub(curdate(), interval ? day)
		group by date(created_at)
		order by date asc
	`, table)
	_ = e.Orm.Raw(sql, days).Scan(&list).Error
	return list
}

func (e EduStats) topResources(column string) []models.EduResource {
	list := make([]models.EduResource, 0)
	_ = e.Orm.Order(column + " desc,id desc").Limit(10).Find(&list).Error
	return list
}

func (e EduStats) topCourses() []models.EduCourse {
	list := make([]models.EduCourse, 0)
	_ = e.Orm.Order("learner_count desc,view_count desc,id desc").Limit(10).Find(&list).Error
	return list
}

func (e EduStats) topActivities() []models.EduActivity {
	list := make([]models.EduActivity, 0)
	_ = e.Orm.Order("signup_count desc,id desc").Limit(10).Find(&list).Error
	return list
}
