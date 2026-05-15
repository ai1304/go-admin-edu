package apis

import (
	"go-admin/app/edu/models"
	"go-admin/common/dto"
	"go-admin/common/objectstorage"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"gorm.io/gorm"
)

type EduResource struct {
	api.Api
}

type resourceQuery struct {
	dto.Pagination
	Keyword          string `form:"keyword"`
	Status           string `form:"status"`
	Sort             string `form:"sort"`
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

type resourceFavoriteReq struct {
	UserId    int    `json:"userId" form:"userId"`
	ClientKey string `json:"clientKey" form:"clientKey"`
}

type resourceCommentReq struct {
	ParentId int    `json:"parentId"`
	UserId   int    `json:"userId"`
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
}

type publicResourceDTO struct {
	models.EduResource
	CoverURL string `json:"coverUrl"`
}

func applyResourceFilters(db *gorm.DB, req resourceQuery) *gorm.DB {
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or keywords like ? or author_name like ?", like, like, like, like)
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
	return db
}

func resourceOrder(sort string) string {
	switch sort {
	case "view":
		return "view_count desc,id desc"
	case "download":
		return "download_count desc,id desc"
	case "favorite":
		return "favorite_count desc,id desc"
	default:
		return "id desc"
	}
}

func (e EduResource) coverURLMap(c *gin.Context, resources []models.EduResource) map[int]string {
	result := make(map[int]string)
	coverIds := make([]int, 0)
	for _, item := range resources {
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
		url, err := storage.PresignedGetObject(c.Request.Context(), file.ObjectKey, 15*time.Minute)
		if err == nil {
			result[file.Id] = url
		}
	}
	return result
}

func (e EduResource) toPublicResources(c *gin.Context, resources []models.EduResource) []publicResourceDTO {
	coverURLs := e.coverURLMap(c, resources)
	list := make([]publicResourceDTO, 0, len(resources))
	for _, item := range resources {
		list = append(list, publicResourceDTO{
			EduResource: item,
			CoverURL:    coverURLs[item.CoverFileId],
		})
	}
	return list
}

func (e EduResource) GetPage(c *gin.Context) {
	req := resourceQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduResource, 0)
	db := applyResourceFilters(e.Orm.Model(&models.EduResource{}), req)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order(resourceOrder(req.Sort)).Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
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
	db := applyResourceFilters(e.Orm.Model(&models.EduResource{}), req)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order(resourceOrder(req.Sort)).Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(e.toPublicResources(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
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
	_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", data.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	data.ViewCount++
	resource := publicResourceDTO{EduResource: data, CoverURL: e.coverURLMap(c, []models.EduResource{data})[data.CoverFileId]}
	e.OK(gin.H{"resource": resource, "files": files}, "查询成功")
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
	_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", data.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	data.ViewCount++
	resource := publicResourceDTO{EduResource: data, CoverURL: e.coverURLMap(c, []models.EduResource{data})[data.CoverFileId]}
	e.OK(gin.H{"resource": resource, "files": files}, "查询成功")
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

func (e EduResource) GetComments(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduResourceComment, 0)
	if err := e.Orm.Where("resource_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduResource) UpdateComment(c *gin.Context) {
	req := models.EduResourceComment{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"nickname":  req.Nickname,
		"content":   req.Content,
		"status":    req.Status,
		"update_by": req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduResourceComment{}).
		Where("id = ? and resource_id = ?", c.Param("commentId"), c.Param("id")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("commentId"), "更新成功")
}

func (e EduResource) DeleteComments(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("resource_id = ?", c.Param("id")).Delete(&models.EduResourceComment{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduResource) PublicFavoriteState(c *gin.Context) {
	req := resourceFavoriteReq{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	if req.UserId == 0 && req.ClientKey == "" {
		e.OK(gin.H{"favorited": false}, "查询成功")
		return
	}
	var count int64
	db := e.Orm.Model(&models.EduResourceFavorite{}).Where("resource_id = ?", c.Param("id"))
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	} else {
		db = db.Where("client_key = ?", req.ClientKey)
	}
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(gin.H{"favorited": count > 0}, "查询成功")
}

func (e EduResource) PublicFavorite(c *gin.Context) {
	req := resourceFavoriteReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	resourceId := parsePathId(c.Param("id"))
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "缺少收藏标识")
		return
	}
	var resource models.EduResource
	if err := e.Orm.Where("id = ? and status = ?", resourceId, models.ResourceStatusPublished).First(&resource).Error; err != nil {
		e.Error(404, err, "资源不存在")
		return
	}
	var count int64
	db := e.Orm.Model(&models.EduResourceFavorite{}).Where("resource_id = ?", resourceId)
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	} else {
		db = db.Where("client_key = ?", req.ClientKey)
	}
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "收藏失败")
		return
	}
	if count > 0 {
		e.OK(resourceId, "已收藏")
		return
	}
	favorite := models.EduResourceFavorite{ResourceId: resourceId, UserId: req.UserId, ClientKey: req.ClientKey}
	if err := e.Orm.Create(&favorite).Error; err != nil {
		e.Error(500, err, "收藏失败")
		return
	}
	_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", resourceId).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	e.OK(resourceId, "收藏成功")
}

func (e EduResource) PublicUnfavorite(c *gin.Context) {
	req := resourceFavoriteReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	resourceId := parsePathId(c.Param("id"))
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "缺少收藏标识")
		return
	}
	db := e.Orm.Where("resource_id = ?", resourceId)
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	} else {
		db = db.Where("client_key = ?", req.ClientKey)
	}
	result := db.Delete(&models.EduResourceFavorite{})
	if result.Error != nil {
		e.Error(500, result.Error, "取消收藏失败")
		return
	}
	if result.RowsAffected > 0 {
		_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", resourceId).UpdateColumn("favorite_count", gorm.Expr("GREATEST(favorite_count - ?, 0)", 1)).Error
	}
	e.OK(resourceId, "取消收藏成功")
}

func (e EduResource) PublicGetComments(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduResourceComment, 0)
	if err := e.Orm.Where("resource_id = ? and status = ?", c.Param("id"), 1).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduResource) PublicCreateComment(c *gin.Context) {
	req := resourceCommentReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.Content == "" {
		e.Error(400, nil, "评论内容不能为空")
		return
	}
	resourceId := parsePathId(c.Param("id"))
	var resource models.EduResource
	if err := e.Orm.Where("id = ? and status = ?", resourceId, models.ResourceStatusPublished).First(&resource).Error; err != nil {
		e.Error(404, err, "资源不存在")
		return
	}
	if req.Nickname == "" {
		req.Nickname = "访客"
	}
	comment := models.EduResourceComment{
		ResourceId: resourceId,
		ParentId:   req.ParentId,
		UserId:     req.UserId,
		Nickname:   req.Nickname,
		Content:    req.Content,
		Status:     1,
	}
	if err := e.Orm.Create(&comment).Error; err != nil {
		e.Error(500, err, "评论失败")
		return
	}
	e.OK(comment, "评论成功")
}
