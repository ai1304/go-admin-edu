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

type EduActivity struct {
	api.Api
}

type activityQuery struct {
	dto.Pagination
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
	SchoolId int    `form:"schoolId"`
}

type activityPortalIdentityReq struct {
	ClientKey string `json:"clientKey" form:"clientKey"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

func (e EduActivity) GetPage(c *gin.Context) {
	req := activityQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduActivity, 0)
	db := applyEduUserScope(c, e.Orm.Model(&models.EduActivity{}))
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
	outcomes := make([]models.EduActivityOutcome, 0)
	_ = e.Orm.Where("activity_id = ? and status = ?", data.Id, 1).Order("id desc").Find(&outcomes).Error
	e.OK(gin.H{"activity": data, "outcomes": outcomes}, "查询成功")
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

func (e EduActivity) PublicSignup(c *gin.Context) {
	req := models.EduActivitySignup{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.ActivityId = parsePathId(c.Param("id"))
	req.UserId = 0
	if req.Status == "" {
		req.Status = "signed"
	}
	if req.ClientKey == "" {
		e.Error(400, nil, "缺少报名标识")
		return
	}
	var activity models.EduActivity
	if err := e.Orm.Where("id = ? and status = ?", req.ActivityId, "published").First(&activity).Error; err != nil {
		e.Error(404, err, "活动不存在")
		return
	}
	var count int64
	if err := e.Orm.Model(&models.EduActivitySignup{}).
		Where("activity_id = ? and client_key = ? and status = ?", req.ActivityId, req.ClientKey, "signed").
		Count(&count).Error; err != nil {
		e.Error(500, err, "报名失败")
		return
	}
	if count > 0 {
		e.OK(req.ActivityId, "already signed")
		return
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "报名失败")
		return
	}
	_ = e.Orm.Model(&models.EduActivity{}).Where("id = ?", req.ActivityId).UpdateColumn("signup_count", gorm.Expr("signup_count + ?", 1)).Error
	e.OK(req.Id, "报名成功")
}

func (e EduActivity) PublicSignupState(c *gin.Context) {
	req := activityPortalIdentityReq{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	if req.ClientKey == "" {
		e.OK(gin.H{"signed": false, "checked": false}, "query success")
		return
	}
	activityId := parsePathId(c.Param("id"))
	var signup models.EduActivitySignup
	signed := e.Orm.Where("activity_id = ? and client_key = ? and status = ?", activityId, req.ClientKey, "signed").First(&signup).Error == nil
	var checkinCount int64
	_ = e.Orm.Model(&models.EduActivityCheckin{}).Where("activity_id = ? and client_key = ? and status = ?", activityId, req.ClientKey, "checked").Count(&checkinCount).Error
	e.OK(gin.H{"signed": signed, "checked": checkinCount > 0, "signup": signup}, "query success")
}

func (e EduActivity) PublicCancelSignup(c *gin.Context) {
	req := activityPortalIdentityReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.ClientKey == "" {
		e.Error(400, nil, "缺少报名标识")
		return
	}
	activityId := parsePathId(c.Param("id"))
	result := e.Orm.Where("activity_id = ? and client_key = ? and status = ?", activityId, req.ClientKey, "signed").
		Delete(&models.EduActivitySignup{})
	if result.Error != nil {
		e.Error(500, result.Error, "取消报名失败")
		return
	}
	if result.RowsAffected > 0 {
		_ = e.Orm.Model(&models.EduActivity{}).Where("id = ?", activityId).UpdateColumn("signup_count", gorm.Expr("GREATEST(signup_count - ?, 0)", 1)).Error
	}
	e.OK(activityId, "取消报名成功")
}

func (e EduActivity) PublicCheckin(c *gin.Context) {
	req := activityPortalIdentityReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.ClientKey == "" {
		e.Error(400, nil, "缺少签到标识")
		return
	}
	activityId := parsePathId(c.Param("id"))
	var activity models.EduActivity
	if err := e.Orm.Where("id = ? and status = ?", activityId, "published").First(&activity).Error; err != nil {
		e.Error(404, err, "活动不存在")
		return
	}
	var signupCount int64
	if err := e.Orm.Model(&models.EduActivitySignup{}).Where("activity_id = ? and client_key = ? and status = ?", activityId, req.ClientKey, "signed").Count(&signupCount).Error; err != nil {
		e.Error(500, err, "签到失败")
		return
	}
	if signupCount == 0 {
		e.Error(400, nil, "请先报名再签到")
		return
	}
	var exists int64
	_ = e.Orm.Model(&models.EduActivityCheckin{}).Where("activity_id = ? and client_key = ? and status = ?", activityId, req.ClientKey, "checked").Count(&exists).Error
	if exists > 0 {
		e.OK(activityId, "already checked")
		return
	}
	checkin := models.EduActivityCheckin{
		ActivityId: activityId,
		ClientKey:  req.ClientKey,
		CheckinAt:  time.Now().Format("2006-01-02 15:04:05"),
		Status:     "checked",
	}
	if err := e.Orm.Create(&checkin).Error; err != nil {
		e.Error(500, err, "签到失败")
		return
	}
	e.OK(checkin, "签到成功")
}

func (e EduActivity) GetSignups(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduActivitySignup, 0)
	if err := e.Orm.Where("activity_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduActivity) InsertSignup(c *gin.Context) {
	req := models.EduActivitySignup{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.ActivityId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == "" {
		req.Status = "signed"
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	_ = e.Orm.Model(&models.EduActivity{}).Where("id = ?", req.ActivityId).UpdateColumn("signup_count", gorm.Expr("signup_count + ?", 1)).Error
	e.OK(req.Id, "创建成功")
}

func (e EduActivity) UpdateSignup(c *gin.Context) {
	req := models.EduActivitySignup{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduActivitySignup{}).
		Where("id = ? and activity_id = ?", c.Param("signupId"), c.Param("id")).
		Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("signupId"), "更新成功")
}

func (e EduActivity) DeleteSignups(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("activity_id = ?", c.Param("id")).Delete(&models.EduActivitySignup{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	_ = e.Orm.Model(&models.EduActivity{}).Where("id = ?", c.Param("id")).UpdateColumn("signup_count", gorm.Expr("GREATEST(signup_count - ?, 0)", len(req.Ids))).Error
	e.OK(req.Ids, "删除成功")
}

func (e EduActivity) GetCheckins(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduActivityCheckin, 0)
	if err := e.Orm.Where("activity_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduActivity) InsertCheckin(c *gin.Context) {
	req := models.EduActivityCheckin{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.ActivityId = parsePathId(c.Param("id"))
	if req.CheckinAt == "" {
		req.CheckinAt = time.Now().Format("2006-01-02 15:04:05")
	}
	if req.Status == "" {
		req.Status = "checked"
	}
	req.SetCreateBy(user.GetUserId(c))
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "签到失败")
		return
	}
	e.OK(req.Id, "签到成功")
}

func (e EduActivity) DeleteCheckins(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("activity_id = ?", c.Param("id")).Delete(&models.EduActivityCheckin{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduActivity) GetOutcomes(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.EduActivityOutcome, 0)
	if err := e.Orm.Where("activity_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

func (e EduActivity) InsertOutcome(c *gin.Context) {
	req := models.EduActivityOutcome{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.ActivityId = parsePathId(c.Param("id"))
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

func (e EduActivity) PublicUploadOutcomeFile(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	activityId := parsePathId(c.Param("id"))
	var activity models.EduActivity
	if err := e.Orm.Where("id = ? and status = ?", activityId, "published").First(&activity).Error; err != nil {
		e.Error(404, err, "activity not found")
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
	objectKey := fmt.Sprintf("tenant/%d/activity/%d/outcome/%s/%s%s", 0, activityId, time.Now().Format("2006/01"), uuid.New().String(), ext)
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
		Usage:        "activity_outcome",
	}
	if err := e.Orm.Create(&record).Error; err != nil {
		e.Error(500, err, "save file record failed")
		return
	}
	e.OK(record, "upload success")
}

func (e EduActivity) PublicSubmitOutcome(c *gin.Context) {
	req := models.EduActivityOutcome{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	activityId := parsePathId(c.Param("id"))
	var activity models.EduActivity
	if err := e.Orm.Where("id = ? and status = ?", activityId, "published").First(&activity).Error; err != nil {
		e.Error(404, err, "activity not found")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		e.Error(400, nil, "outcome title is required")
		return
	}
	req.ActivityId = activityId
	if req.Status == 0 {
		req.Status = 1
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "submit outcome failed")
		return
	}
	e.OK(req, "submit success")
}

func (e EduActivity) UpdateOutcome(c *gin.Context) {
	req := models.EduActivityOutcome{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduActivityOutcome{}).
		Where("id = ? and activity_id = ?", c.Param("outcomeId"), c.Param("id")).
		Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("outcomeId"), "更新成功")
}

func (e EduActivity) DeleteOutcomes(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Where("activity_id = ?", c.Param("id")).Delete(&models.EduActivityOutcome{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
