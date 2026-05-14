package apis

import (
	"fmt"
	"go-admin/app/edu/models"
	"go-admin/common/dto"
	"go-admin/common/objectstorage"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EduResourceFile struct {
	api.Api
}

type fileQuery struct {
	dto.Pagination
	ResourceId int    `form:"resourceId"`
	Usage      string `form:"usage"`
}

func (e EduResourceFile) GetPage(c *gin.Context) {
	req := fileQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduResourceFile, 0)
	db := e.Orm.Model(&models.EduResourceFile{})
	if req.ResourceId != 0 {
		db = db.Where("resource_id = ?", req.ResourceId)
	}
	if req.Usage != "" {
		db = db.Where("usage = ?", req.Usage)
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

func (e EduResourceFile) Insert(c *gin.Context) {
	req := models.EduResourceFile{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetCreateBy(user.GetUserId(c))
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduResourceFile) Upload(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	multipartFile, err := c.FormFile("file")
	if err != nil {
		e.Error(400, err, "文件不能为空")
		return
	}
	file, err := multipartFile.Open()
	if err != nil {
		e.Error(500, err, "读取文件失败")
		return
	}
	defer file.Close()

	storage, err := objectstorage.NewFromExtend()
	if err != nil {
		e.Error(500, err, "对象存储未配置")
		return
	}
	if err := storage.EnsureBucket(c.Request.Context()); err != nil {
		e.Error(500, err, "初始化存储桶失败")
		return
	}

	resourceId, _ := strconv.Atoi(c.PostForm("resourceId"))
	usage := c.DefaultPostForm("usage", "attachment")
	ext := strings.ToLower(filepath.Ext(multipartFile.Filename))
	objectKey := fmt.Sprintf("tenant/%d/resource/%s/%s%s", 0, time.Now().Format("2006/01"), uuid.New().String(), ext)
	contentType := multipartFile.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if err := storage.PutObject(c.Request.Context(), objectKey, file, multipartFile.Size, contentType); err != nil {
		e.Error(500, err, "文件上传失败")
		return
	}

	record := models.EduResourceFile{
		ResourceId:   resourceId,
		OriginalName: multipartFile.Filename,
		ObjectKey:    objectKey,
		Bucket:       storage.BucketName(),
		ContentType:  contentType,
		Ext:          strings.TrimPrefix(ext, "."),
		Size:         multipartFile.Size,
		Usage:        usage,
	}
	record.SetCreateBy(user.GetUserId(c))
	if err := e.Orm.Create(&record).Error; err != nil {
		e.Error(500, err, "文件记录保存失败")
		return
	}
	e.OK(record, "上传成功")
}

func (e EduResourceFile) PublicAccessURL(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}

	resourceId := parsePathId(c.Param("id"))
	fileId := parsePathId(c.Param("fileId"))
	var resource models.EduResource
	if err := e.Orm.Where("id = ? and status = ?", resourceId, models.ResourceStatusPublished).First(&resource).Error; err != nil {
		e.Error(404, err, "资源不存在")
		return
	}

	var file models.EduResourceFile
	if err := e.Orm.Where("id = ? and resource_id = ?", fileId, resourceId).First(&file).Error; err != nil {
		e.Error(404, err, "文件不存在")
		return
	}

	storage, err := objectstorage.NewFromExtend()
	if err != nil {
		e.Error(500, err, "对象存储未配置")
		return
	}
	expires := 15 * time.Minute
	url, err := storage.PresignedGetObject(c.Request.Context(), file.ObjectKey, expires)
	if err != nil {
		e.Error(500, err, "获取文件访问地址失败")
		return
	}

	_ = e.Orm.Model(&models.EduResource{}).Where("id = ?", resource.Id).UpdateColumn("download_count", gorm.Expr("download_count + ?", 1)).Error
	e.OK(gin.H{
		"url":       url,
		"expiresIn": int(expires.Seconds()),
		"file":      file,
	}, "获取成功")
}

func (e EduResourceFile) Delete(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Delete(&models.EduResourceFile{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
