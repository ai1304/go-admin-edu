package apis

import (
	"fmt"
	"go-admin/app/edu/models"
	"go-admin/common/dto"
	"go-admin/common/objectstorage"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EduCourse struct {
	api.Api
}

type courseQuery struct {
	dto.Pagination
	Keyword          string `form:"keyword"`
	Status           string `form:"status"`
	SchoolId         int    `form:"schoolId"`
	StageCategoryId  int    `form:"stageCategoryId"`
	DisabilityTypeId int    `form:"disabilityTypeId"`
	Category         string `form:"category"`
	Difficulty       string `form:"difficulty"`
	Sort             string `form:"sort"`
}

type publicCourseIdentityReq struct {
	UserId    int    `form:"userId" json:"userId"`
	ClientKey string `form:"clientKey" json:"clientKey"`
}

type publicLearningReq struct {
	UserId         int    `json:"userId"`
	ClientKey      string `json:"clientKey"`
	Progress       int    `json:"progress"`
	WatchedSeconds int    `json:"watchedSeconds"`
	Status         string `json:"status"`
}

type publicAssignmentSubmissionReq struct {
	UserId    int    `json:"userId"`
	ClientKey string `json:"clientKey"`
	Nickname  string `json:"nickname"`
	Content   string `json:"content"`
	FileId    int    `json:"fileId"`
}

type publicCourseDTO struct {
	models.EduCourse
	CoverURL string `json:"coverUrl"`
}

func (e EduCourse) courseCoverURLMap(c *gin.Context, courses []models.EduCourse) map[int]string {
	result := make(map[int]string)
	coverIds := make([]int, 0)
	for _, item := range courses {
		if item.CoverFileId != 0 {
			coverIds = append(coverIds, item.CoverFileId)
		}
	}
	if len(coverIds) == 0 {
		return result
	}
	files := make([]models.EduResourceFile, 0)
	if err := e.Orm.Where("id in ?", coverIds).Find(&files).Error; err != nil {
		return result
	}
	storage, err := objectstorage.NewFromExtend()
	if err != nil {
		return result
	}
	for _, file := range files {
		if url, err := storage.PresignedGetObject(c.Request.Context(), file.ObjectKey, 15*time.Minute); err == nil {
			result[file.Id] = url
		}
	}
	return result
}

func (e EduCourse) toPublicCourses(c *gin.Context, courses []models.EduCourse) []publicCourseDTO {
	coverURLs := e.courseCoverURLMap(c, courses)
	result := make([]publicCourseDTO, 0, len(courses))
	for _, item := range courses {
		result = append(result, publicCourseDTO{EduCourse: item, CoverURL: coverURLs[item.CoverFileId]})
	}
	return result
}

func (e EduCourse) GetPage(c *gin.Context) {
	req := courseQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCourse, 0)
	db := applyEduUserScope(c, e.Orm.Model(&models.EduCourse{}))
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or teacher_name like ?", like, like, like)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.SchoolId != 0 {
		db = db.Where("school_id = ?", req.SchoolId)
	}
	if req.StageCategoryId != 0 {
		db = db.Where("stage_category_id = ?", req.StageCategoryId)
	}
	if req.DisabilityTypeId != 0 {
		db = db.Where("disability_type_id = ?", req.DisabilityTypeId)
	}
	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Difficulty != "" {
		db = db.Where("difficulty = ?", req.Difficulty)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order("sort desc,id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(e.toPublicCourses(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e EduCourse) PublicGetPage(c *gin.Context) {
	req := courseQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCourse, 0)
	db := e.Orm.Model(&models.EduCourse{}).Where("status = ?", "published")
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or teacher_name like ?", like, like, like)
	}
	if req.StageCategoryId != 0 {
		db = db.Where("stage_category_id = ?", req.StageCategoryId)
	}
	if req.DisabilityTypeId != 0 {
		db = db.Where("disability_type_id = ?", req.DisabilityTypeId)
	}
	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Difficulty != "" {
		db = db.Where("difficulty = ?", req.Difficulty)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order(courseOrder(req.Sort)).Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func courseOrder(sort string) string {
	switch sort {
	case "view":
		return "view_count desc,id desc"
	case "learner":
		return "learner_count desc,id desc"
	default:
		return "sort desc,id desc"
	}
}

func (e EduCourse) Get(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduCourse
	if err := e.Orm.First(&data, c.Param("id")).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	chapters := make([]models.EduCourseChapter, 0)
	lessons := make([]models.EduCourseLesson, 0)
	_ = e.Orm.Where("course_id = ?", data.Id).Order("sort asc,id asc").Find(&chapters).Error
	_ = e.Orm.Where("course_id = ?", data.Id).Order("sort asc,id asc").Find(&lessons).Error
	e.OK(gin.H{"course": data, "chapters": chapters, "lessons": lessons}, "查询成功")
}

func (e EduCourse) PublicGet(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduCourse
	if err := e.Orm.Where("status = ?", "published").First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "课程不存在")
		return
	}
	chapters := make([]models.EduCourseChapter, 0)
	lessons := make([]models.EduCourseLesson, 0)
	assignments := make([]models.EduAssignment, 0)
	_ = e.Orm.Where("course_id = ? and status = ?", data.Id, 1).Order("sort asc,id asc").Find(&chapters).Error
	_ = e.Orm.Where("course_id = ? and status = ?", data.Id, 1).Order("sort asc,id asc").Find(&lessons).Error
	_ = e.Orm.Where("course_id = ? and status = ?", data.Id, 1).Order("id desc").Find(&assignments).Error
	_ = e.Orm.Model(&models.EduCourse{}).Where("id = ?", data.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	data.ViewCount++
	videoURL := ""
	if data.VideoFileId != 0 {
		if url, err := e.courseVideoURL(c, data.VideoFileId); err == nil {
			videoURL = url
		}
	}
	course := e.toPublicCourses(c, []models.EduCourse{data})[0]
	e.OK(gin.H{"course": course, "chapters": chapters, "lessons": lessons, "assignments": assignments, "videoUrl": videoURL}, "查询成功")
}

func (e EduCourse) courseVideoURL(c *gin.Context, fileId int) (string, error) {
	var file models.EduResourceFile
	if err := e.Orm.First(&file, fileId).Error; err != nil {
		return "", err
	}
	storage, err := objectstorage.NewFromExtend()
	if err != nil {
		return "", err
	}
	return storage.PresignedGetObject(c.Request.Context(), file.ObjectKey, 15*time.Minute)
}

func (e EduCourse) Insert(c *gin.Context) {
	req := models.EduCourse{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == "" {
		req.Status = "draft"
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCourse) Update(c *gin.Context) {
	req := models.EduCourse{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduCourse{}).Where("id = ?", c.Param("id")).Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("id"), "更新成功")
}

func (e EduCourse) Delete(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Delete(&models.EduCourse{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCourse) PublicGetAssignments(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	courseId := parsePathId(c.Param("id"))
	var course models.EduCourse
	if err := e.Orm.Where("id = ? and status = ?", courseId, "published").First(&course).Error; err != nil {
		e.Error(404, err, "课程不存在")
		return
	}
	list := make([]models.EduAssignment, 0)
	if err := e.Orm.Where("course_id = ? and status = ?", courseId, 1).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduCourse) PublicGetLearningRecords(c *gin.Context) {
	req := publicCourseIdentityReq{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	if req.UserId == 0 && req.ClientKey == "" {
		e.OK([]models.EduLearningRecord{}, "查询成功")
		return
	}
	list := make([]models.EduLearningRecord, 0)
	db := e.Orm.Where("course_id = ?", c.Param("id"))
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	} else {
		db = db.Where("client_key = ?", req.ClientKey)
	}
	if err := db.Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduCourse) PublicLessonVideoURL(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	courseId := parsePathId(c.Param("id"))
	lessonId := parsePathId(c.Param("lessonId"))
	var course models.EduCourse
	if err := e.Orm.Where("id = ? and status = ?", courseId, "published").First(&course).Error; err != nil {
		e.Error(404, err, "course not found")
		return
	}
	var lesson models.EduCourseLesson
	if err := e.Orm.Where("id = ? and course_id = ? and status = ?", lessonId, courseId, 1).First(&lesson).Error; err != nil {
		e.Error(404, err, "lesson not found")
		return
	}
	if lesson.VideoFileId == 0 {
		e.Error(400, nil, "lesson video is empty")
		return
	}
	var file models.EduResourceFile
	if err := e.Orm.First(&file, lesson.VideoFileId).Error; err != nil {
		e.Error(404, err, "video file not found")
		return
	}
	storage, err := objectstorage.NewFromExtend()
	if err != nil {
		e.Error(500, err, "storage is not configured")
		return
	}
	expires := 15 * time.Minute
	url, err := storage.PresignedGetObject(c.Request.Context(), file.ObjectKey, expires)
	if err != nil {
		e.Error(500, err, "get video url failed")
		return
	}
	e.OK(gin.H{"url": url, "expiresIn": int(expires.Seconds()), "file": file, "lesson": lesson}, "query success")
}

func (e EduCourse) PublicTrackLearning(c *gin.Context) {
	req := publicLearningReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	courseId := parsePathId(c.Param("id"))
	lessonId := parsePathId(c.Param("lessonId"))
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "缺少学习标识")
		return
	}
	if req.Progress < 0 {
		req.Progress = 0
	}
	if req.Progress > 100 {
		req.Progress = 100
	}
	if req.WatchedSeconds < 0 {
		req.WatchedSeconds = 0
	}
	if req.Status == "" {
		req.Status = "learning"
	}
	if req.Progress >= 100 {
		req.Status = "finished"
	}
	var course models.EduCourse
	if err := e.Orm.Where("id = ? and status = ?", courseId, "published").First(&course).Error; err != nil {
		e.Error(404, err, "课程不存在")
		return
	}
	var lesson models.EduCourseLesson
	if err := e.Orm.Where("id = ? and course_id = ? and status = ?", lessonId, courseId, 1).First(&lesson).Error; err != nil {
		e.Error(404, err, "课时不存在")
		return
	}
	identity := e.Orm.Model(&models.EduLearningRecord{}).Where("course_id = ?", courseId)
	if req.UserId != 0 {
		identity = identity.Where("user_id = ?", req.UserId)
	} else {
		identity = identity.Where("client_key = ?", req.ClientKey)
	}
	var identityCount int64
	_ = identity.Count(&identityCount).Error

	var record models.EduLearningRecord
	db := e.Orm.Where("course_id = ? and lesson_id = ?", courseId, lessonId)
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	} else {
		db = db.Where("client_key = ?", req.ClientKey)
	}
	err := db.First(&record).Error
	if err == nil {
		if err := e.Orm.Model(&record).Updates(map[string]interface{}{
			"progress":        req.Progress,
			"watched_seconds": req.WatchedSeconds,
			"status":          req.Status,
		}).Error; err != nil {
			e.Error(500, err, "记录学习失败")
			return
		}
		record.Progress = req.Progress
		record.WatchedSeconds = req.WatchedSeconds
		record.Status = req.Status
		e.OK(record, "记录成功")
		return
	}
	if err != gorm.ErrRecordNotFound {
		e.Error(500, err, "记录学习失败")
		return
	}
	record = models.EduLearningRecord{
		CourseId:       courseId,
		LessonId:       lessonId,
		UserId:         req.UserId,
		ClientKey:      req.ClientKey,
		Progress:       req.Progress,
		WatchedSeconds: req.WatchedSeconds,
		Status:         req.Status,
	}
	if err := e.Orm.Create(&record).Error; err != nil {
		e.Error(500, err, "记录学习失败")
		return
	}
	if identityCount == 0 {
		_ = e.Orm.Model(&models.EduCourse{}).Where("id = ?", courseId).UpdateColumn("learner_count", gorm.Expr("learner_count + ?", 1)).Error
	}
	e.OK(record, "记录成功")
}

func (e EduCourse) PublicSubmitAssignment(c *gin.Context) {
	req := publicAssignmentSubmissionReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	courseId := parsePathId(c.Param("id"))
	assignmentId := parsePathId(c.Param("assignmentId"))
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "缺少提交标识")
		return
	}
	if req.Content == "" && req.FileId == 0 {
		e.Error(400, nil, "提交内容不能为空")
		return
	}
	var course models.EduCourse
	if err := e.Orm.Where("id = ? and status = ?", courseId, "published").First(&course).Error; err != nil {
		e.Error(404, err, "课程不存在")
		return
	}
	var assignment models.EduAssignment
	if err := e.Orm.Where("id = ? and course_id = ? and status = ?", assignmentId, courseId, 1).First(&assignment).Error; err != nil {
		e.Error(404, err, "作业不存在")
		return
	}
	if req.Nickname == "" {
		req.Nickname = "学习者"
	}
	submission := models.EduAssignmentSubmission{
		CourseId:     courseId,
		AssignmentId: assignmentId,
		UserId:       req.UserId,
		ClientKey:    req.ClientKey,
		Nickname:     req.Nickname,
		Content:      req.Content,
		FileId:       req.FileId,
		Status:       "submitted",
	}
	if err := e.Orm.Create(&submission).Error; err != nil {
		e.Error(500, err, "提交失败")
		return
	}
	e.OK(submission, "提交成功")
}

func (e EduCourse) PublicUploadAssignmentFile(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	courseId := parsePathId(c.Param("id"))
	assignmentId := parsePathId(c.Param("assignmentId"))
	var course models.EduCourse
	if err := e.Orm.Where("id = ? and status = ?", courseId, "published").First(&course).Error; err != nil {
		e.Error(404, err, "course not found")
		return
	}
	var assignment models.EduAssignment
	if err := e.Orm.Where("id = ? and course_id = ? and status = ?", assignmentId, courseId, 1).First(&assignment).Error; err != nil {
		e.Error(404, err, "assignment not found")
		return
	}
	multipartFile, err := c.FormFile("file")
	if err != nil {
		e.Error(400, err, "file is required")
		return
	}
	file, err := multipartFile.Open()
	if err != nil {
		e.Error(500, err, "read file failed")
		return
	}
	defer file.Close()

	storage, err := objectstorage.NewFromExtend()
	if err != nil {
		e.Error(500, err, "storage is not configured")
		return
	}
	if err := storage.EnsureBucket(c.Request.Context()); err != nil {
		e.Error(500, err, "init storage bucket failed")
		return
	}
	ext := strings.ToLower(filepath.Ext(multipartFile.Filename))
	objectKey := fmt.Sprintf("tenant/%d/course/%d/assignment/%s/%s%s", 0, courseId, time.Now().Format("2006/01"), uuid.New().String(), ext)
	contentType := multipartFile.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if err := storage.PutObject(c.Request.Context(), objectKey, file, multipartFile.Size, contentType); err != nil {
		e.Error(500, err, "upload file failed")
		return
	}
	record := models.EduResourceFile{
		OriginalName: multipartFile.Filename,
		ObjectKey:    objectKey,
		Bucket:       storage.BucketName(),
		ContentType:  contentType,
		Ext:          strings.TrimPrefix(ext, "."),
		Size:         multipartFile.Size,
		Usage:        "assignment",
	}
	if err := e.Orm.Create(&record).Error; err != nil {
		e.Error(500, err, "save file record failed")
		return
	}
	e.OK(record, "upload success")
}

func (e EduCourse) GetChapters(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduCourseChapter, 0)
	if err := e.Orm.Where("course_id = ?", c.Param("id")).Order("sort asc,id asc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduCourse) InsertChapter(c *gin.Context) {
	req := models.EduCourseChapter{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.CourseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == 0 {
		req.Status = 1
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCourse) UpdateChapter(c *gin.Context) {
	req := models.EduCourseChapter{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduCourseChapter{}).
		Where("id = ? and course_id = ?", c.Param("chapterId"), c.Param("id")).
		Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("chapterId"), "更新成功")
}

func (e EduCourse) DeleteChapters(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("course_id = ?", c.Param("id")).Delete(&models.EduCourseChapter{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCourse) GetLessons(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduCourseLesson, 0)
	db := e.Orm.Where("course_id = ?", c.Param("id"))
	if chapterId := c.Query("chapterId"); chapterId != "" {
		db = db.Where("chapter_id = ?", chapterId)
	}
	if err := db.Order("sort asc,id asc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduCourse) InsertLesson(c *gin.Context) {
	req := models.EduCourseLesson{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.CourseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == 0 {
		req.Status = 1
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCourse) UpdateLesson(c *gin.Context) {
	req := models.EduCourseLesson{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduCourseLesson{}).
		Where("id = ? and course_id = ?", c.Param("lessonId"), c.Param("id")).
		Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("lessonId"), "更新成功")
}

func (e EduCourse) DeleteLessons(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("course_id = ?", c.Param("id")).Delete(&models.EduCourseLesson{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCourse) GetAssignments(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduAssignment, 0)
	if err := e.Orm.Where("course_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduCourse) InsertAssignment(c *gin.Context) {
	req := models.EduAssignment{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.CourseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == 0 {
		req.Status = 1
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCourse) UpdateAssignment(c *gin.Context) {
	req := models.EduAssignment{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"title":     req.Title,
		"content":   req.Content,
		"status":    req.Status,
		"update_by": req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduAssignment{}).
		Where("id = ? and course_id = ?", c.Param("assignmentId"), c.Param("id")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("assignmentId"), "更新成功")
}

func (e EduCourse) DeleteAssignments(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("course_id = ?", c.Param("id")).Delete(&models.EduAssignment{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCourse) GetAssignmentSubmissions(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduAssignmentSubmission, 0)
	if err := e.Orm.Where("course_id = ? and assignment_id = ?", c.Param("id"), c.Param("assignmentId")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduCourse) InsertAssignmentSubmission(c *gin.Context) {
	req := models.EduAssignmentSubmission{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.CourseId = parsePathId(c.Param("id"))
	req.AssignmentId = parsePathId(c.Param("assignmentId"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == "" {
		req.Status = "submitted"
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCourse) UpdateAssignmentSubmission(c *gin.Context) {
	req := models.EduAssignmentSubmission{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"user_id":   req.UserId,
		"content":   req.Content,
		"file_id":   req.FileId,
		"score":     req.Score,
		"status":    req.Status,
		"update_by": req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduAssignmentSubmission{}).
		Where("id = ? and course_id = ? and assignment_id = ?", c.Param("submissionId"), c.Param("id"), c.Param("assignmentId")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("submissionId"), "更新成功")
}

func (e EduCourse) GetAssignmentSubmissionFileURL(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var submission models.EduAssignmentSubmission
	if err := e.Orm.Where("id = ? and course_id = ? and assignment_id = ?", c.Param("submissionId"), c.Param("id"), c.Param("assignmentId")).
		First(&submission).Error; err != nil {
		e.Error(404, err, "submission not found")
		return
	}
	if submission.FileId == 0 {
		e.Error(400, nil, "submission file is empty")
		return
	}
	var file models.EduResourceFile
	if err := e.Orm.First(&file, submission.FileId).Error; err != nil {
		e.Error(404, err, "file not found")
		return
	}
	storage, err := objectstorage.NewFromExtend()
	if err != nil {
		e.Error(500, err, "storage is not configured")
		return
	}
	expires := 15 * time.Minute
	url, err := storage.PresignedGetObject(c.Request.Context(), file.ObjectKey, expires)
	if err != nil {
		e.Error(500, err, "get file url failed")
		return
	}
	e.OK(gin.H{"url": url, "expiresIn": int(expires.Seconds()), "file": file}, "query success")
}

func (e EduCourse) DeleteAssignmentSubmissions(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("course_id = ? and assignment_id = ?", c.Param("id"), c.Param("assignmentId")).
		Delete(&models.EduAssignmentSubmission{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCourse) GetLearningRecords(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduLearningRecord, 0)
	db := e.Orm.Where("course_id = ?", c.Param("id"))
	if lessonId := c.Query("lessonId"); lessonId != "" {
		db = db.Where("lesson_id = ?", lessonId)
	}
	if err := db.Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduCourse) InsertLearningRecord(c *gin.Context) {
	req := models.EduLearningRecord{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.CourseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == "" {
		req.Status = "learning"
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCourse) UpdateLearningRecord(c *gin.Context) {
	req := models.EduLearningRecord{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"lesson_id":       req.LessonId,
		"user_id":         req.UserId,
		"progress":        req.Progress,
		"watched_seconds": req.WatchedSeconds,
		"status":          req.Status,
		"update_by":       req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduLearningRecord{}).
		Where("id = ? and course_id = ?", c.Param("recordId"), c.Param("id")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("recordId"), "更新成功")
}

func (e EduCourse) DeleteLearningRecords(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("course_id = ?", c.Param("id")).Delete(&models.EduLearningRecord{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
