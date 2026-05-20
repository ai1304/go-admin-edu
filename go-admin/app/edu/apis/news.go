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

type EduNews struct {
	api.Api
}

type newsQuery struct {
	dto.Pagination
	Keyword    string `form:"keyword"`
	ModuleType string `form:"moduleType"`
	Status     string `form:"status"`
}

type publicNewsDTO struct {
	models.EduNews
	CoverURL string `json:"coverUrl"`
}

func applyNewsFilters(db *gorm.DB, req newsQuery) *gorm.DB {
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or source like ? or keywords like ?", like, like, like, like)
	}
	if req.ModuleType != "" {
		db = db.Where("module_type = ?", req.ModuleType)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	return db
}

func (e EduNews) newsCoverURLMap(c *gin.Context, list []models.EduNews) map[int]string {
	result := make(map[int]string)
	ids := make([]int, 0)
	for _, item := range list {
		if item.CoverFileId != 0 {
			ids = append(ids, item.CoverFileId)
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

func (e EduNews) toPublicNews(c *gin.Context, list []models.EduNews) []publicNewsDTO {
	coverURLs := e.newsCoverURLMap(c, list)
	result := make([]publicNewsDTO, 0, len(list))
	for _, item := range list {
		cover := item.CoverUrl
		if item.CoverFileId != 0 && coverURLs[item.CoverFileId] != "" {
			cover = coverURLs[item.CoverFileId]
		}
		result = append(result, publicNewsDTO{EduNews: item, CoverURL: cover})
	}
	return result
}

func (e EduNews) GetPage(c *gin.Context) {
	req := newsQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduNews, 0)
	db := applyEduUserScope(c, applyNewsFilters(e.Orm.Model(&models.EduNews{}), req))
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	if err := db.Order("is_top desc,sort desc,publish_time desc,id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.PageOK(e.toPublicNews(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "query success")
}

func (e EduNews) PublicGetPage(c *gin.Context) {
	req := newsQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	req.Status = models.NewsStatusPublished
	list := make([]models.EduNews, 0)
	db := applyNewsFilters(e.Orm.Model(&models.EduNews{}), req)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	if err := db.Order("is_top desc,sort desc,publish_time desc,id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.PageOK(e.toPublicNews(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "query success")
}

func (e EduNews) Get(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduNews
	if err := e.Orm.First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "news not found")
		return
	}
	e.OK(e.toPublicNews(c, []models.EduNews{data})[0], "query success")
}

func (e EduNews) PublicGet(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduNews
	if err := e.Orm.Where("status = ?", models.NewsStatusPublished).First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "news not found")
		return
	}
	_ = e.Orm.Model(&models.EduNews{}).Where("id = ?", data.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	data.ViewCount++
	e.OK(e.toPublicNews(c, []models.EduNews{data})[0], "query success")
}

func (e EduNews) Insert(c *gin.Context) {
	req := models.EduNews{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.Status == "" {
		req.Status = models.NewsStatusDraft
	}
	req.SetCreateBy(user.GetUserId(c))
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "create failed")
		return
	}
	e.OK(req.Id, "create success")
}

func (e EduNews) Update(c *gin.Context) {
	req := models.EduNews{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduNews{}).Where("id = ?", c.Param("id")).Updates(&req).Error; err != nil {
		e.Error(500, err, "update failed")
		return
	}
	e.OK(c.Param("id"), "update success")
}

func (e EduNews) Delete(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Delete(&models.EduNews{}, req.Ids).Error; err != nil {
		e.Error(500, err, "delete failed")
		return
	}
	e.OK(req.Ids, "delete success")
}

func (e EduNews) PublicLike(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Model(&models.EduNews{}).Where("id = ? and status = ?", c.Param("id"), models.NewsStatusPublished).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		e.Error(500, err, "like failed")
		return
	}
	e.OK(c.Param("id"), "like success")
}
