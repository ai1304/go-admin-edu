package apis

import (
	"go-admin/app/edu/models"
	"go-admin/common/dto"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

type EduResource struct {
	api.Api
}

type resourceQuery struct {
	dto.Pagination
	Keyword          string `form:"keyword"`
	Status           string `form:"status"`
	SchoolId         int    `form:"schoolId"`
	StageCategoryId  int    `form:"stageCategoryId"`
	DisabilityTypeId int    `form:"disabilityTypeId"`
	ResourceTypeId   int    `form:"resourceTypeId"`
	AbilityDomainId  int    `form:"abilityDomainId"`
	TopicCategoryId  int    `form:"topicCategoryId"`
}

type resourceReviewReq struct {
	Action  string `json:"action"`
	Comment string `json:"comment"`
}

func (e EduResource) GetPage(c *gin.Context) {
	req := resourceQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduResource, 0)
	db := e.Orm.Model(&models.EduResource{})
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or keywords like ?", like, like, like)
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
	if req.ResourceTypeId != 0 {
		db = db.Where("resource_type_id = ?", req.ResourceTypeId)
	}
	if req.AbilityDomainId != 0 {
		db = db.Where("ability_domain_id = ?", req.AbilityDomainId)
	}
	if req.TopicCategoryId != 0 {
		db = db.Where("topic_category_id = ?", req.TopicCategoryId)
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

func (e EduResource) PublicGetPage(c *gin.Context) {
	req := resourceQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	req.Status = models.ResourceStatusPublished
	list := make([]models.EduResource, 0)
	db := e.Orm.Model(&models.EduResource{}).Where("status = ?", req.Status)
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or keywords like ?", like, like, like)
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

func (e EduResource) Get(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduResource
	if err := e.Orm.First(&data, c.Param("id")).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	files := make([]models.EduResourceFile, 0)
	_ = e.Orm.Where("resource_id = ?", data.Id).Find(&files).Error
	e.OK(gin.H{"resource": data, "files": files}, "查询成功")
}

func (e EduResource) PublicGet(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduResource
	if err := e.Orm.Where("status = ?", models.ResourceStatusPublished).First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "资源不存在")
		return
	}
	files := make([]models.EduResourceFile, 0)
	_ = e.Orm.Where("resource_id = ?", data.Id).Find(&files).Error
	e.OK(gin.H{"resource": data, "files": files}, "查询成功")
}

func (e EduResource) Insert(c *gin.Context) {
	req := models.EduResource{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == "" {
		req.Status = models.ResourceStatusDraft
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduResource) Update(c *gin.Context) {
	req := models.EduResource{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduResource{}).Where("id = ?", c.Param("id")).Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("id"), "更新成功")
}

func (e EduResource) Delete(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Delete(&models.EduResource{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduResource) SubmitReview(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Model(&models.EduResource{}).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{
		"status":    models.ResourceStatusReviewing,
		"update_by": user.GetUserId(c),
	}).Error; err != nil {
		e.Error(500, err, "提交审核失败")
		return
	}
	e.OK(c.Param("id"), "提交审核成功")
}

func (e EduResource) Review(c *gin.Context) {
	req := resourceReviewReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var resource models.EduResource
	if err := e.Orm.First(&resource, c.Param("id")).Error; err != nil {
		e.Error(500, err, "资源不存在")
		return
	}
	afterStatus := models.ResourceStatusRejected
	if req.Action == "approve" {
		afterStatus = models.ResourceStatusPublished
	}
	review := models.EduResourceReview{
		ResourceId:   resource.Id,
		Action:       req.Action,
		Comment:      req.Comment,
		BeforeStatus: resource.Status,
		AfterStatus:  afterStatus,
	}
	review.SetCreateBy(user.GetUserId(c))
	tx := e.Orm.Begin()
	if err := tx.Model(&models.EduResource{}).Where("id = ?", resource.Id).Updates(map[string]interface{}{
		"status":    afterStatus,
		"update_by": user.GetUserId(c),
	}).Error; err != nil {
		tx.Rollback()
		e.Error(500, err, "审核失败")
		return
	}
	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		e.Error(500, err, "审核失败")
		return
	}
	tx.Commit()
	e.OK(resource.Id, "审核成功")
}
