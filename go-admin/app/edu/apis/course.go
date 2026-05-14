package apis

import (
	"go-admin/app/edu/models"
	"go-admin/common/dto"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
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
}

func (e EduCourse) GetPage(c *gin.Context) {
	req := courseQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCourse, 0)
	db := e.Orm.Model(&models.EduCourse{})
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
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order("id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
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
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order("id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
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
	_ = e.Orm.Where("course_id = ? and status = ?", data.Id, 1).Order("sort asc,id asc").Find(&chapters).Error
	_ = e.Orm.Where("course_id = ? and status = ?", data.Id, 1).Order("sort asc,id asc").Find(&lessons).Error
	e.OK(gin.H{"course": data, "chapters": chapters, "lessons": lessons}, "查询成功")
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
