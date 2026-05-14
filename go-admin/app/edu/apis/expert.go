package apis

import (
	"go-admin/app/edu/models"
	"go-admin/common/dto"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

type EduExpert struct {
	api.Api
}

type expertQuery struct {
	dto.Pagination
	Keyword string `form:"keyword"`
	Status  int    `form:"status"`
}

func (e EduExpert) GetPage(c *gin.Context) {
	req := expertQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduExpert, 0)
	db := e.Orm.Model(&models.EduExpert{})
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("name like ? or title like ? or specialties like ?", like, like, like)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
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

func (e EduExpert) PublicGetPage(c *gin.Context) {
	req := expertQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduExpert, 0)
	db := e.Orm.Model(&models.EduExpert{}).Where("status = ?", 1)
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("name like ? or title like ? or specialties like ?", like, like, like)
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

func (e EduExpert) Get(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduExpert
	if err := e.Orm.First(&data, c.Param("id")).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	resources := make([]models.EduExpertResource, 0)
	_ = e.Orm.Where("expert_id = ?", data.Id).Order("id desc").Find(&resources).Error
	e.OK(gin.H{"expert": data, "resources": resources}, "查询成功")
}

func (e EduExpert) PublicGet(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduExpert
	if err := e.Orm.Where("status = ?", 1).First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "专家不存在")
		return
	}
	e.OK(data, "查询成功")
}

func (e EduExpert) Insert(c *gin.Context) {
	req := models.EduExpert{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
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

func (e EduExpert) Update(c *gin.Context) {
	req := models.EduExpert{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduExpert{}).Where("id = ?", c.Param("id")).Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("id"), "更新成功")
}

func (e EduExpert) Delete(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Delete(&models.EduExpert{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
