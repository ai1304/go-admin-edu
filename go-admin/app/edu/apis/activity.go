package apis

import (
	"go-admin/app/edu/models"
	"go-admin/common/dto"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"gorm.io/gorm"
)

type EduActivity struct {
	api.Api
}

type activityQuery struct {
	dto.Pagination
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
	SchoolId int    `form:"schoolId"`
}

func (e EduActivity) GetPage(c *gin.Context) {
	req := activityQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduActivity, 0)
	db := e.Orm.Model(&models.EduActivity{})
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or organizer like ?", like, like, like)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.SchoolId != 0 {
		db = db.Where("school_id = ?", req.SchoolId)
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

func (e EduActivity) PublicGetPage(c *gin.Context) {
	req := activityQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduActivity, 0)
	db := e.Orm.Model(&models.EduActivity{}).Where("status = ?", "published")
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or organizer like ?", like, like, like)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order("start_time desc,id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e EduActivity) Get(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduActivity
	if err := e.Orm.First(&data, c.Param("id")).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	outcomes := make([]models.EduActivityOutcome, 0)
	_ = e.Orm.Where("activity_id = ?", data.Id).Order("id desc").Find(&outcomes).Error
	e.OK(gin.H{"activity": data, "outcomes": outcomes}, "查询成功")
}

func (e EduActivity) PublicGet(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduActivity
	if err := e.Orm.Where("status = ?", "published").First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "活动不存在")
		return
	}
	e.OK(data, "查询成功")
}

func (e EduActivity) Insert(c *gin.Context) {
	req := models.EduActivity{}
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

func (e EduActivity) Update(c *gin.Context) {
	req := models.EduActivity{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduActivity{}).Where("id = ?", c.Param("id")).Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("id"), "更新成功")
}

func (e EduActivity) Delete(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Delete(&models.EduActivity{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduActivity) Signup(c *gin.Context) {
	req := models.EduActivitySignup{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.ActivityId = parsePathId(c.Param("id"))
	req.UserId = user.GetUserId(c)
	req.SetCreateBy(req.UserId)
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "报名失败")
		return
	}
	_ = e.Orm.Model(&models.EduActivity{}).Where("id = ?", req.ActivityId).UpdateColumn("signup_count", gorm.Expr("signup_count + ?", 1)).Error
	e.OK(req.Id, "报名成功")
}
