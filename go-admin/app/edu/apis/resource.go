package apis

import (
	"go-admin/app/edu/models"
	edusearch "go-admin/app/edu/search"
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
	TagId            int    `form:"tagId"`
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

type resourceStatusReq struct {
	Status  string `json:"status"`
	Comment string `json:"comment"`
}

type resourceSaveReq struct {
	models.EduResource
	TagIds []int `json:"tagIds"`
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
	TagIds   []int                   `json:"tagIds"`
	Tags     []models.EduResourceTag `json:"tags"`
	CoverURL string                  `json:"coverUrl"`
}

func toSearchResourceQuery(req resourceQuery) edusearch.ResourceQuery {
	return edusearch.ResourceQuery{
		Keyword:          req.Keyword,
		Status:           req.Status,
		Sort:             req.Sort,
		TagId:            req.TagId,
		SchoolId:         req.SchoolId,
		StageCategoryId:  req.StageCategoryId,
		DisabilityTypeId: req.DisabilityTypeId,
		ResourceTypeId:   req.ResourceTypeId,
		AbilityDomainId:  req.AbilityDomainId,
		TopicCategoryId:  req.TopicCategoryId,
		PageIndex:        req.GetPageIndex(),
		PageSize:         req.GetPageSize(),
	}
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
	if req.TagId != 0 {
		db = db.Where("id in (?)", db.Session(&gorm.Session{}).
			Model(&models.EduResourceTagRelation{}).
			Select("resource_id").
			Where("tag_id = ?", req.TagId))
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
	tagMap := e.tagMapForResources(resources)
	list := make([]publicResourceDTO, 0, len(resources))
	for _, item := range resources {
		tags := tagMap[item.Id]
		list = append(list, publicResourceDTO{
			EduResource: item,
			TagIds:      tagIdsFromTags(tags),
			Tags:        tags,
			CoverURL:    coverURLs[item.CoverFileId],
		})
	}
	return list
}

func tagIdsFromTags(tags []models.EduResourceTag) []int {
	ids := make([]int, 0, len(tags))
	for _, item := range tags {
		ids = append(ids, item.Id)
	}
	return ids
}

func (e EduResource) tagMapForResources(resources []models.EduResource) map[int][]models.EduResourceTag {
	result := make(map[int][]models.EduResourceTag)
	resourceIds := make([]int, 0, len(resources))
	for _, item := range resources {
		resourceIds = append(resourceIds, item.Id)
	}
	if len(resourceIds) == 0 {
		return result
	}
	relations := make([]models.EduResourceTagRelation, 0)
	if err := e.Orm.Where("resource_id in ?", resourceIds).Find(&relations).Error; err != nil {
		return result
	}
	tagIds := make([]int, 0, len(relations))
	for _, item := range relations {
		tagIds = append(tagIds, item.TagId)
	}
	if len(tagIds) == 0 {
		return result
	}
	tags := make([]models.EduResourceTag, 0)
	if err := e.Orm.Where("id in ? and status = ?", tagIds, 1).Find(&tags).Error; err != nil {
		return result
	}
	tagById := make(map[int]models.EduResourceTag)
	for _, item := range tags {
		tagById[item.Id] = item
	}
	for _, relation := range relations {
		if tag, ok := tagById[relation.TagId]; ok {
			result[relation.ResourceId] = append(result[relation.ResourceId], tag)
		}
	}
	return result
}

func (e EduResource) syncResourceTags(resourceId int, tagIds []int, currentUserId int) error {
	if err := e.Orm.Where("resource_id = ?", resourceId).Delete(&models.EduResourceTagRelation{}).Error; err != nil {
		return err
	}
	for _, tagId := range tagIds {
		if tagId == 0 {
			continue
		}
		relation := models.EduResourceTagRelation{ResourceId: resourceId, TagId: tagId}
		relation.SetCreateBy(currentUserId)
		if err := e.Orm.Create(&relation).Error; err != nil {
			return err
		}
	}
	return nil
}

func (e EduResource) GetPage(c *gin.Context) {
	req := resourceQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduResource, 0)
	db := applyEduUserScope(c, applyResourceFilters(e.Orm.Model(&models.EduResource{}), req))
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	if err := db.Order(resourceOrder(req.Sort)).Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.PageOK(e.toPublicResources(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "query success")
}

func (e EduResource) PublicGetPage(c *gin.Context) {
	req := resourceQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	req.Status = models.ResourceStatusPublished
	searcher := edusearch.NewMySQLSearcher(e.Orm)
	result, err := searcher.SearchResources(toSearchResourceQuery(req))
	if err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.PageOK(e.toPublicResources(c, result.List), int(result.Count), req.GetPageIndex(), req.GetPageSize(), "query success")
}

func (e EduResource) PublicSearch(c *gin.Context) {
	e.PublicGetPage(c)
}

func (e EduResource) Get(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduResource
	if err := e.Orm.First(&data, c.Param("id")).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	files := make([]models.EduResourceFile, 0)
	_ = e.Orm.Where("resource_id = ?", data.Id).Find(&files).Error
	_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", data.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	data.ViewCount++
	resource := e.toPublicResources(c, []models.EduResource{data})[0]
	e.OK(gin.H{"resource": resource, "files": files}, "query success")
}

func (e EduResource) PublicGet(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduResource
	if err := e.Orm.Where("status = ?", models.ResourceStatusPublished).First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "resource not found")
		return
	}
	files := make([]models.EduResourceFile, 0)
	_ = e.Orm.Where("resource_id = ?", data.Id).Find(&files).Error
	_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", data.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	data.ViewCount++
	resource := e.toPublicResources(c, []models.EduResource{data})[0]
	e.OK(gin.H{"resource": resource, "files": files}, "query success")
}

func (e EduResource) Insert(c *gin.Context) {
	req := resourceSaveReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	resource := req.EduResource
	currentUserId := user.GetUserId(c)
	resource.SetCreateBy(currentUserId)
	if resource.Status == "" {
		resource.Status = models.ResourceStatusDraft
	}
	if err := e.Orm.Create(&resource).Error; err != nil {
		e.Error(500, err, "create failed")
		return
	}
	if err := e.syncResourceTags(resource.Id, req.TagIds, currentUserId); err != nil {
		e.Error(500, err, "save tags failed")
		return
	}
	e.OK(resource.Id, "create success")
}

func (e EduResource) Update(c *gin.Context) {
	req := resourceSaveReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	resource := req.EduResource
	currentUserId := user.GetUserId(c)
	resource.SetUpdateBy(currentUserId)
	if err := e.Orm.Model(&models.EduResource{}).Where("id = ?", c.Param("id")).Updates(&resource).Error; err != nil {
		e.Error(500, err, "update failed")
		return
	}
	if err := e.syncResourceTags(parsePathId(c.Param("id")), req.TagIds, currentUserId); err != nil {
		e.Error(500, err, "save tags failed")
		return
	}
	e.OK(c.Param("id"), "update success")
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
		e.Error(500, err, "delete failed")
		return
	}
	_ = e.Orm.Where("resource_id in ?", req.Ids).Delete(&models.EduResourceTagRelation{}).Error
	e.OK(req.Ids, "delete success")
}

func (e EduResource) ReindexSearch(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var total int64
	var published int64
	if err := e.Orm.Model(&models.EduResource{}).Count(&total).Error; err != nil {
		e.Error(500, err, "count resource failed")
		return
	}
	if err := e.Orm.Model(&models.EduResource{}).Where("status = ?", models.ResourceStatusPublished).Count(&published).Error; err != nil {
		e.Error(500, err, "count resource failed")
		return
	}
	e.OK(gin.H{
		"engine":    "mysql",
		"total":     total,
		"published": published,
		"synced":    published,
	}, "search index synced")
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
		e.Error(500, err, "submit review failed")
		return
	}
	e.OK(c.Param("id"), "submit review success")
}

func (e EduResource) Review(c *gin.Context) {
	req := resourceReviewReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var resource models.EduResource
	if err := e.Orm.First(&resource, c.Param("id")).Error; err != nil {
		e.Error(500, err, "resource not found")
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
		e.Error(500, err, "review failed")
		return
	}
	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		e.Error(500, err, "review failed")
		return
	}
	tx.Commit()
	e.OK(resource.Id, "review success")
}

func (e EduResource) GetReviews(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduResourceReview, 0)
	if err := e.Orm.Where("resource_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.OK(list, "query success")
}

func (e EduResource) UpdateStatus(c *gin.Context) {
	req := resourceStatusReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.Status != models.ResourceStatusPublished && req.Status != models.ResourceStatusOffline {
		e.Error(400, nil, "unsupported status")
		return
	}
	var resource models.EduResource
	if err := e.Orm.First(&resource, c.Param("id")).Error; err != nil {
		e.Error(404, err, "resource not found")
		return
	}
	if resource.Status == req.Status {
		e.OK(resource.Id, "status unchanged")
		return
	}
	if resource.Status != models.ResourceStatusPublished && resource.Status != models.ResourceStatusOffline {
		e.Error(400, nil, "only published or offline resource can change status")
		return
	}
	action := "publish"
	if req.Status == models.ResourceStatusOffline {
		action = "offline"
	}
	review := models.EduResourceReview{
		ResourceId:   resource.Id,
		Action:       action,
		Comment:      req.Comment,
		BeforeStatus: resource.Status,
		AfterStatus:  req.Status,
	}
	review.SetCreateBy(user.GetUserId(c))
	tx := e.Orm.Begin()
	if err := tx.Model(&models.EduResource{}).Where("id = ?", resource.Id).Updates(map[string]interface{}{
		"status":    req.Status,
		"update_by": user.GetUserId(c),
	}).Error; err != nil {
		tx.Rollback()
		e.Error(500, err, "status update failed")
		return
	}
	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		e.Error(500, err, "status update failed")
		return
	}
	tx.Commit()
	e.OK(resource.Id, "status update success")
}

func (e EduResource) GetComments(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduResourceComment, 0)
	if err := e.Orm.Where("resource_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.OK(list, "query success")
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
		e.Error(500, err, "update failed")
		return
	}
	e.OK(c.Param("commentId"), "update success")
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
		e.Error(500, err, "delete failed")
		return
	}
	e.OK(req.Ids, "delete success")
}

func (e EduResource) PublicFavoriteState(c *gin.Context) {
	req := resourceFavoriteReq{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	if req.UserId == 0 && req.ClientKey == "" {
		e.OK(gin.H{"favorited": false}, "query success")
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
		e.Error(500, err, "query failed")
		return
	}
	e.OK(gin.H{"favorited": count > 0}, "query success")
}

func (e EduResource) PublicFavorite(c *gin.Context) {
	req := resourceFavoriteReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	resourceId := parsePathId(c.Param("id"))
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "missing favorite identity")
		return
	}
	var resource models.EduResource
	if err := e.Orm.Where("id = ? and status = ?", resourceId, models.ResourceStatusPublished).First(&resource).Error; err != nil {
		e.Error(404, err, "resource not found")
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
		e.Error(500, err, "favorite failed")
		return
	}
	if count > 0 {
		e.OK(resourceId, "already favorited")
		return
	}
	favorite := models.EduResourceFavorite{ResourceId: resourceId, UserId: req.UserId, ClientKey: req.ClientKey}
	if err := e.Orm.Create(&favorite).Error; err != nil {
		e.Error(500, err, "favorite failed")
		return
	}
	_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", resourceId).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	e.OK(resourceId, "favorite success")
}

func (e EduResource) PublicUnfavorite(c *gin.Context) {
	req := resourceFavoriteReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	resourceId := parsePathId(c.Param("id"))
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "missing favorite identity")
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
		e.Error(500, result.Error, "unfavorite failed")
		return
	}
	if result.RowsAffected > 0 {
		_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", resourceId).UpdateColumn("favorite_count", gorm.Expr("GREATEST(favorite_count - ?, 0)", 1)).Error
	}
	e.OK(resourceId, "unfavorite success")
}

func (e EduResource) PublicGetComments(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduResourceComment, 0)
	if err := e.Orm.Where("resource_id = ? and status = ?", c.Param("id"), 1).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.OK(list, "query success")
}

func (e EduResource) PublicCreateComment(c *gin.Context) {
	req := resourceCommentReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.Content == "" {
		e.Error(400, nil, "comment content required")
		return
	}
	resourceId := parsePathId(c.Param("id"))
	var resource models.EduResource
	if err := e.Orm.Where("id = ? and status = ?", resourceId, models.ResourceStatusPublished).First(&resource).Error; err != nil {
		e.Error(404, err, "resource not found")
		return
	}
	if req.ParentId != 0 {
		var parentCount int64
		if err := e.Orm.Model(&models.EduResourceComment{}).
			Where("id = ? and resource_id = ? and status = ?", req.ParentId, resourceId, 1).
			Count(&parentCount).Error; err != nil {
			e.Error(500, err, "comment failed")
			return
		}
		if parentCount == 0 {
			e.Error(400, nil, "parent comment not found")
			return
		}
	}
	if req.Nickname == "" {
		req.Nickname = "guest"
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
		e.Error(500, err, "comment failed")
		return
	}
	e.OK(comment, "comment success")
}

func (e EduResource) PublicLikeComment(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	resourceId := parsePathId(c.Param("id"))
	commentId := parsePathId(c.Param("commentId"))
	result := e.Orm.Model(&models.EduResourceComment{}).
		Where("id = ? and resource_id = ? and status = ?", commentId, resourceId, 1).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
	if result.Error != nil {
		e.Error(500, result.Error, "like failed")
		return
	}
	if result.RowsAffected == 0 {
		e.Error(404, nil, "comment not found")
		return
	}
	e.OK(commentId, "like success")
}
