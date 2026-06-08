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

type EduExpert struct {
	api.Api
}

type expertQuery struct {
	dto.Pagination
	Keyword       string `form:"keyword"`
	Title         string `form:"title"`
	Specialty     string `form:"specialty"`
	IsRecommended int    `form:"isRecommended"`
	Status        int    `form:"status"`
	Sort          string `form:"sort"`
}

type expertFavoriteReq struct {
	UserId    int    `json:"userId" form:"userId"`
	ClientKey string `json:"clientKey" form:"clientKey"`
}

type publicExpertDTO struct {
	models.EduExpert
	AvatarURL string `json:"avatarUrl"`
}

func (e EduExpert) expertAvatarURLMap(c *gin.Context, list []models.EduExpert) map[int]string {
	result := make(map[int]string)
	ids := make([]int, 0)
	for _, item := range list {
		if item.AvatarFileId != 0 {
			ids = append(ids, item.AvatarFileId)
		}
	}
	if len(ids) == 0 {
		return result
	}
	files := make([]models.EduResourceFile, 0)
	if err := e.Orm.Where("id in ?", ids).Find(&files).Error; err != nil {
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

func (e EduExpert) toPublicExperts(c *gin.Context, list []models.EduExpert) []publicExpertDTO {
	avatarURLs := e.expertAvatarURLMap(c, list)
	result := make([]publicExpertDTO, 0, len(list))
	for _, item := range list {
		result = append(result, publicExpertDTO{EduExpert: item, AvatarURL: avatarURLs[item.AvatarFileId]})
	}
	return result
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
		db = db.Where("name like ? or title like ? or organization like ? or specialties like ?", like, like, like, like)
	}
	if req.Title != "" {
		db = db.Where("title = ?", req.Title)
	}
	if req.Specialty != "" {
		db = db.Where("specialties like ?", "%"+req.Specialty+"%")
	}
	if req.IsRecommended != 0 {
		db = db.Where("is_recommended = ?", req.IsRecommended)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order(expertOrder(req.Sort)).Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(e.toPublicExperts(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func expertOrder(sort string) string {
	switch sort {
	case "view":
		return "view_count desc,id desc"
	default:
		return "is_recommended desc,sort desc,id desc"
	}
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
		db = db.Where("name like ? or title like ? or organization like ? or specialties like ?", like, like, like, like)
	}
	if req.Title != "" {
		db = db.Where("title = ?", req.Title)
	}
	if req.Specialty != "" {
		db = db.Where("specialties like ?", "%"+req.Specialty+"%")
	}
	if req.IsRecommended != 0 {
		db = db.Where("is_recommended = ?", req.IsRecommended)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order("is_recommended desc,sort desc,id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(e.toPublicExperts(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
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
	_ = e.Orm.Model(&models.EduExpert{}).Where("id = ?", data.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	data.ViewCount++
	resources := make([]models.EduExpertResource, 0)
	_ = e.Orm.Where("expert_id = ? and status = ?", data.Id, 1).Order("id desc").Find(&resources).Error
	rankings := make([]models.EduExpert, 0)
	_ = e.Orm.Where("status = ?", 1).Order("view_count desc,id desc").Limit(10).Find(&rankings).Error
	e.OK(gin.H{"expert": e.toPublicExperts(c, []models.EduExpert{data})[0], "resources": resources, "rankings": e.toPublicExperts(c, rankings)}, "查询成功")
}

func (e EduExpert) PublicResourceAccessURL(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	expertId := parsePathId(c.Param("id"))
	resourceId := parsePathId(c.Param("resourceId"))
	var relation models.EduExpertResource
	if err := e.Orm.Where("id = ? and expert_id = ? and type = ? and status = ?", resourceId, expertId, "file", 1).First(&relation).Error; err != nil {
		e.Error(404, err, "expert resource not found")
		return
	}
	if relation.FileId == 0 {
		e.Error(400, nil, "expert resource file is empty")
		return
	}
	var file models.EduResourceFile
	if err := e.Orm.First(&file, relation.FileId).Error; err != nil {
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

func (e EduExpert) PublicFavoriteState(c *gin.Context) {
	req := expertFavoriteReq{}
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
	db := e.Orm.Model(&models.EduExpertFavorite{}).Where("expert_id = ?", c.Param("id"))
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

func (e EduExpert) PublicFavorite(c *gin.Context) {
	req := expertFavoriteReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	expertId := parsePathId(c.Param("id"))
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "missing favorite identity")
		return
	}
	var expert models.EduExpert
	if err := e.Orm.Where("id = ? and status = ?", expertId, 1).First(&expert).Error; err != nil {
		e.Error(404, err, "expert not found")
		return
	}
	var count int64
	db := e.Orm.Model(&models.EduExpertFavorite{}).Where("expert_id = ?", expertId)
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
		e.OK(expertId, "already favorited")
		return
	}
	favorite := models.EduExpertFavorite{ExpertId: expertId, UserId: req.UserId, ClientKey: req.ClientKey}
	if err := e.Orm.Create(&favorite).Error; err != nil {
		e.Error(500, err, "favorite failed")
		return
	}
	_ = e.Orm.Model(&models.EduExpert{}).Where("id = ?", expertId).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	e.OK(expertId, "favorite success")
}

func (e EduExpert) PublicUnfavorite(c *gin.Context) {
	req := expertFavoriteReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	expertId := parsePathId(c.Param("id"))
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "missing favorite identity")
		return
	}
	db := e.Orm.Where("expert_id = ?", expertId)
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	} else {
		db = db.Where("client_key = ?", req.ClientKey)
	}
	result := db.Delete(&models.EduExpertFavorite{})
	if result.Error != nil {
		e.Error(500, result.Error, "unfavorite failed")
		return
	}
	if result.RowsAffected > 0 {
		_ = e.Orm.Model(&models.EduExpert{}).Where("id = ?", expertId).UpdateColumn("favorite_count", gorm.Expr("GREATEST(favorite_count - ?, 0)", 1)).Error
	}
	e.OK(expertId, "unfavorite success")
}

func (e EduExpert) PublicShare(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	expertId := parsePathId(c.Param("id"))
	result := e.Orm.Model(&models.EduExpert{}).
		Where("id = ? and status = ?", expertId, 1).
		UpdateColumn("share_count", gorm.Expr("share_count + ?", 1))
	if result.Error != nil {
		e.Error(500, result.Error, "share failed")
		return
	}
	if result.RowsAffected == 0 {
		e.Error(404, nil, "expert not found")
		return
	}
	e.OK(expertId, "share success")
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

func (e EduExpert) GetResources(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduExpertResource, 0)
	if err := e.Orm.Where("expert_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduExpert) InsertResource(c *gin.Context) {
	req := models.EduExpertResource{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.ExpertId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if err := e.Orm.Select("ExpertId", "Title", "Summary", "Type", "ResourceId", "CourseId", "FileId", "Status", "CreateBy").Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduExpert) UpdateResource(c *gin.Context) {
	req := models.EduExpertResource{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"title":       req.Title,
		"summary":     req.Summary,
		"type":        req.Type,
		"resource_id": req.ResourceId,
		"course_id":   req.CourseId,
		"file_id":     req.FileId,
		"status":      req.Status,
		"update_by":   req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduExpertResource{}).
		Where("id = ? and expert_id = ?", c.Param("resourceId"), c.Param("id")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("resourceId"), "更新成功")
}

func (e EduExpert) DeleteResources(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("expert_id = ?", c.Param("id")).Delete(&models.EduExpertResource{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
