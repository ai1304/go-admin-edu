package apis

import (
	"encoding/csv"
	"fmt"
	"go-admin/app/edu/models"
	"go-admin/common/dto"
	"go-admin/common/objectstorage"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"gorm.io/gorm"
)

type EduCase struct {
	api.Api
}

type caseQuery struct {
	dto.Pagination
	Keyword        string `form:"keyword"`
	Status         string `form:"status"`
	SchoolId       int    `form:"schoolId"`
	Stage          string `form:"stage"`
	DisabilityType string `form:"disabilityType"`
	AbilityDomain  string `form:"abilityDomain"`
	CaseType       string `form:"caseType"`
	Desensitize    bool   `form:"desensitize"`
}

type caseAccessLogQuery struct {
	dto.Pagination
	Action  string `form:"action"`
	UserId  int    `form:"userId"`
	Keyword string `form:"keyword"`
}

type caseAuthorizationQuery struct {
	dto.Pagination
	UserId int    `form:"userId"`
	Scope  string `form:"scope"`
	Status string `form:"status"`
}

type caseReviewReq struct {
	Action  string `json:"action"`
	Comment string `json:"comment"`
}

func (e EduCase) caseCoverURLMap(c *gin.Context, list []models.EduCase) map[int]string {
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

func (e EduCase) fillCaseCoverURLs(c *gin.Context, list []models.EduCase) []models.EduCase {
	coverURLs := e.caseCoverURLMap(c, list)
	for index := range list {
		if coverURL := coverURLs[list[index].CoverFileId]; coverURL != "" {
			list[index].CoverUrl = coverURL
		}
	}
	return list
}

func applyCaseAccessLogFilters(db *gorm.DB, req caseAccessLogQuery) *gorm.DB {
	if req.Action != "" {
		db = db.Where("action = ?", req.Action)
	}
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	}
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("ip like ? or path like ? or user_agent like ?", like, like, like)
	}
	return db
}

func (e EduCase) writeAccessLog(c *gin.Context, caseId int, action string) {
	log := models.EduCaseAccessLog{
		CaseId:    caseId,
		UserId:    user.GetUserId(c),
		Action:    action,
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
		Ip:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	log.SetCreateBy(log.UserId)
	_ = e.Orm.Create(&log).Error
}

func desensitizeCases(list []models.EduCase) {
	for index := range list {
		desensitizeCase(&list[index])
	}
}

func desensitizeCase(data *models.EduCase) {
	data.StudentName = maskName(data.StudentName)
	data.StudentCode = maskCode(data.StudentCode)
	data.Birthday = ""
	data.Summary = maskLongText(data.Summary)
}

func desensitizeIEPs(list []models.EduCaseIEP) {
	for index := range list {
		list[index].Goal = maskLongText(list[index].Goal)
		list[index].Plan = maskLongText(list[index].Plan)
		list[index].Evaluation = maskLongText(list[index].Evaluation)
	}
}

func desensitizeAssessments(list []models.EduCaseAssessment) {
	for index := range list {
		list[index].Result = maskLongText(list[index].Result)
	}
}

func desensitizeInterventions(list []models.EduCaseIntervention) {
	for index := range list {
		list[index].Content = maskLongText(list[index].Content)
	}
}

func maskName(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) == 1 {
		return "*"
	}
	return string(runes[0]) + "**"
}

func maskCode(value string) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= 4 {
		return "****"
	}
	return string(runes[:2]) + "****" + string(runes[len(runes)-2:])
}

func maskLongText(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return "内容已脱敏"
}

func shouldDesensitize(c *gin.Context) bool {
	value := strings.ToLower(c.Query("desensitize"))
	return value == "true" || value == "1" || value == "yes"
}

func allowedAuthorizationScopes(required string) []string {
	switch required {
	case "review":
		return []string{"review"}
	case "edit":
		return []string{"edit", "review"}
	default:
		return []string{"view", "edit", "review"}
	}
}

func authorizationDate(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 10 {
		return value[:10]
	}
	return value
}

func isAuthorizationEffective(data models.EduCaseAuthorization) bool {
	if data.Status != "active" {
		return false
	}
	today := time.Now().Format("2006-01-02")
	startAt := authorizationDate(data.StartAt)
	endAt := authorizationDate(data.EndAt)
	if startAt != "" && startAt > today {
		return false
	}
	if endAt != "" && endAt < today {
		return false
	}
	return true
}

func (e EduCase) hasCaseAccess(c *gin.Context, caseData models.EduCase, requiredScope string) bool {
	currentUserId := user.GetUserId(c)
	if currentUserId == 1 || caseData.CreateBy == currentUserId {
		return true
	}
	list := make([]models.EduCaseAuthorization, 0)
	err := e.Orm.Where("case_id = ? and user_id = ? and scope in ?", caseData.Id, currentUserId, allowedAuthorizationScopes(requiredScope)).
		Find(&list).Error
	if err != nil {
		return false
	}
	for _, item := range list {
		if isAuthorizationEffective(item) {
			return true
		}
	}
	return false
}

func (e EduCase) ensureCaseAccess(c *gin.Context, caseId int, requiredScope string, deniedAction string) (models.EduCase, bool) {
	var data models.EduCase
	if err := e.Orm.First(&data, caseId).Error; err != nil {
		e.Error(500, err, "查询失败")
		return data, false
	}
	if !e.hasCaseAccess(c, data, requiredScope) {
		e.writeAccessLog(c, caseId, deniedAction)
		e.Error(403, nil, "无权访问该案例")
		return data, false
	}
	return data, true
}

func (e EduCase) GetPage(c *gin.Context) {
	req := caseQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCase, 0)
	db := applyEduUserScope(c, e.Orm.Model(&models.EduCase{}))
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or student_name like ? or student_code like ?", like, like, like)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.SchoolId != 0 {
		db = db.Where("school_id = ?", req.SchoolId)
	}
	if req.Stage != "" {
		db = db.Where("stage = ?", req.Stage)
	}
	if req.DisabilityType != "" {
		db = db.Where("disability_type = ?", req.DisabilityType)
	}
	if req.AbilityDomain != "" {
		db = db.Where("ability_domain = ?", req.AbilityDomain)
	}
	if req.CaseType != "" {
		db = db.Where("case_type = ?", req.CaseType)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order("sort desc,id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if req.Desensitize {
		desensitizeCases(list)
	}
	e.PageOK(e.fillCaseCoverURLs(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e EduCase) PublicGetPage(c *gin.Context) {
	req := caseQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCase, 0)
	db := e.Orm.Model(&models.EduCase{}).Where("status = ?", "published")
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("title like ? or summary like ? or school like ? or disability_type like ?", like, like, like, like)
	}
	if req.Stage != "" {
		db = db.Where("stage = ?", req.Stage)
	}
	if req.DisabilityType != "" {
		db = db.Where("disability_type = ?", req.DisabilityType)
	}
	if req.AbilityDomain != "" {
		db = db.Where("ability_domain = ?", req.AbilityDomain)
	}
	if req.CaseType != "" {
		db = db.Where("case_type = ?", req.CaseType)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order("sort desc,id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	desensitizeCases(list)
	e.PageOK(e.fillCaseCoverURLs(c, list), int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e EduCase) PublicGet(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduCase
	if err := e.Orm.Where("status = ?", "published").First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "case not found")
		return
	}
	ieps := make([]models.EduCaseIEP, 0)
	assessments := make([]models.EduCaseAssessment, 0)
	interventions := make([]models.EduCaseIntervention, 0)
	_ = e.Orm.Where("case_id = ?", data.Id).Order("id desc").Find(&ieps).Error
	_ = e.Orm.Where("case_id = ?", data.Id).Order("id desc").Find(&assessments).Error
	_ = e.Orm.Where("case_id = ?", data.Id).Order("id desc").Find(&interventions).Error
	_ = e.Orm.Model(&models.EduCase{}).Where("id = ?", data.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	data.ViewCount++
	desensitizeCase(&data)
	desensitizeIEPs(ieps)
	desensitizeAssessments(assessments)
	desensitizeInterventions(interventions)
	data = e.fillCaseCoverURLs(c, []models.EduCase{data})[0]
	e.OK(gin.H{"case": data, "ieps": ieps, "assessments": assessments, "interventions": interventions}, "查询成功")
}

func (e EduCase) Get(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	data, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "view", "view_detail_denied")
	if !ok {
		return
	}
	ieps := make([]models.EduCaseIEP, 0)
	assessments := make([]models.EduCaseAssessment, 0)
	interventions := make([]models.EduCaseIntervention, 0)
	_ = e.Orm.Where("case_id = ?", data.Id).Order("id desc").Find(&ieps).Error
	_ = e.Orm.Where("case_id = ?", data.Id).Order("id desc").Find(&assessments).Error
	_ = e.Orm.Where("case_id = ?", data.Id).Order("id desc").Find(&interventions).Error
	if shouldDesensitize(c) {
		desensitizeCase(&data)
		desensitizeIEPs(ieps)
		desensitizeAssessments(assessments)
		desensitizeInterventions(interventions)
	}
	data = e.fillCaseCoverURLs(c, []models.EduCase{data})[0]
	e.writeAccessLog(c, data.Id, "view_detail")
	e.OK(gin.H{"case": data, "ieps": ieps, "assessments": assessments, "interventions": interventions}, "查询成功")
}

func (e EduCase) Insert(c *gin.Context) {
	req := models.EduCase{}
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

func (e EduCase) Update(c *gin.Context) {
	req := models.EduCase{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "edit_case_denied"); !ok {
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduCase{}).Where("id = ?", c.Param("id")).Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("id"), "更新成功")
}

func (e EduCase) Delete(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	for _, id := range req.Ids {
		if _, ok := e.ensureCaseAccess(c, id, "edit", "delete_case_denied"); !ok {
			return
		}
	}
	if err := e.Orm.Delete(&models.EduCase{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCase) SubmitReview(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	data, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "submit_case_review_denied")
	if !ok {
		return
	}
	beforeStatus := data.Status
	if err := e.Orm.Model(&models.EduCase{}).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"status":    "reviewing",
		"update_by": user.GetUserId(c),
	}).Error; err != nil {
		e.Error(500, err, "submit review failed")
		return
	}
	review := models.EduCaseReview{CaseId: data.Id, Action: "submit", BeforeStatus: beforeStatus, AfterStatus: "reviewing"}
	review.SetCreateBy(user.GetUserId(c))
	_ = e.Orm.Create(&review).Error
	e.OK(data.Id, "submit review success")
}

func (e EduCase) Review(c *gin.Context) {
	req := caseReviewReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	data, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "review", "review_case_denied")
	if !ok {
		return
	}
	afterStatus := "rejected"
	if req.Action == "approve" {
		afterStatus = "archived"
	}
	review := models.EduCaseReview{
		CaseId:       data.Id,
		Action:       req.Action,
		Comment:      req.Comment,
		BeforeStatus: data.Status,
		AfterStatus:  afterStatus,
	}
	review.SetCreateBy(user.GetUserId(c))
	tx := e.Orm.Begin()
	if err := tx.Model(&models.EduCase{}).Where("id = ?", data.Id).Updates(map[string]interface{}{
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
	e.OK(data.Id, "review success")
}

func (e EduCase) GetReviews(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "review", "view_case_reviews_denied"); !ok {
		return
	}
	list := make([]models.EduCaseReview, 0)
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.OK(list, "query success")
}

func (e EduCase) GetAttachments(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "view", "view_case_attachments_denied"); !ok {
		return
	}
	list := make([]models.EduCaseAttachment, 0)
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.OK(list, "query success")
}

func (e EduCase) AddAttachment(c *gin.Context) {
	req := models.EduCaseAttachment{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "add_case_attachment_denied"); !ok {
		return
	}
	req.CaseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == 0 {
		req.Status = 1
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "create failed")
		return
	}
	e.OK(req.Id, "create success")
}

func (e EduCase) DeleteAttachments(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "delete_case_attachments_denied"); !ok {
		return
	}
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Delete(&models.EduCaseAttachment{}, req.Ids).Error; err != nil {
		e.Error(500, err, "delete failed")
		return
	}
	e.OK(req.Ids, "delete success")
}

func (e EduCase) GetAttachmentURL(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "view", "view_case_attachment_file_denied"); !ok {
		return
	}
	var attachment models.EduCaseAttachment
	if err := e.Orm.Where("id = ? and case_id = ?", c.Param("attachmentId"), c.Param("id")).First(&attachment).Error; err != nil {
		e.Error(404, err, "attachment not found")
		return
	}
	var file models.EduResourceFile
	if err := e.Orm.First(&file, attachment.FileId).Error; err != nil {
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

func (e EduCase) GetAccessLogs(c *gin.Context) {
	req := caseAccessLogQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "review", "view_access_logs_denied"); !ok {
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCaseAccessLog, 0)
	db := applyCaseAccessLogFilters(e.Orm.Model(&models.EduCaseAccessLog{}).Where("case_id = ?", c.Param("id")), req)
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

func (e EduCase) ExportAccessLogs(c *gin.Context) {
	req := caseAccessLogQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "review", "export_access_logs_denied"); !ok {
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCaseAccessLog, 0)
	db := applyCaseAccessLogFilters(e.Orm.Model(&models.EduCaseAccessLog{}).Where("case_id = ?", c.Param("id")), req)
	if err := db.Order("id desc").Limit(10000).Find(&list).Error; err != nil {
		e.Error(500, err, "export access logs failed")
		return
	}

	filename := fmt.Sprintf("case-%s-access-logs-%s.csv", c.Param("id"), time.Now().Format("20060102150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	_ = writer.Write([]string{"ID", "Case ID", "User ID", "Action", "Method", "Path", "IP", "User-Agent", "Created At"})
	for _, item := range list {
		_ = writer.Write([]string{
			strconv.Itoa(item.Id),
			strconv.Itoa(item.CaseId),
			strconv.Itoa(item.UserId),
			item.Action,
			item.Method,
			item.Path,
			item.Ip,
			item.UserAgent,
			item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}

func (e EduCase) GetAuthorizations(c *gin.Context) {
	req := caseAuthorizationQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "review", "view_authorizations_denied"); !ok {
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCaseAuthorization, 0)
	db := e.Orm.Model(&models.EduCaseAuthorization{}).Where("case_id = ?", c.Param("id"))
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	}
	if req.Scope != "" {
		db = db.Where("scope = ?", req.Scope)
	}
	if req.Status != "" {
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

func (e EduCase) AddAuthorization(c *gin.Context) {
	req := models.EduCaseAuthorization{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "review", "add_authorization_denied"); !ok {
		return
	}
	req.CaseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Scope == "" {
		req.Scope = "view"
	}
	if req.Status == "" {
		req.Status = "active"
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建授权失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCase) UpdateAuthorization(c *gin.Context) {
	req := models.EduCaseAuthorization{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "review", "update_authorization_denied"); !ok {
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"user_id":   req.UserId,
		"scope":     req.Scope,
		"start_at":  req.StartAt,
		"end_at":    req.EndAt,
		"status":    req.Status,
		"remark":    req.Remark,
		"update_by": req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduCaseAuthorization{}).
		Where("id = ? and case_id = ?", c.Param("authorizationId"), c.Param("id")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新授权失败")
		return
	}
	e.OK(c.Param("authorizationId"), "更新成功")
}

func (e EduCase) DeleteAuthorizations(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "review", "delete_authorization_denied"); !ok {
		return
	}
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Delete(&models.EduCaseAuthorization{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除授权失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCase) AddIEP(c *gin.Context) {
	req := models.EduCaseIEP{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "add_iep_denied"); !ok {
		return
	}
	req.CaseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == "" {
		req.Status = "draft"
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建IEP失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCase) GetIEPs(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "view", "view_ieps_denied"); !ok {
		return
	}
	list := make([]models.EduCaseIEP, 0)
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if shouldDesensitize(c) {
		desensitizeIEPs(list)
	}
	e.writeAccessLog(c, parsePathId(c.Param("id")), "view_ieps")
	e.OK(list, "查询成功")
}

func (e EduCase) UpdateIEP(c *gin.Context) {
	req := models.EduCaseIEP{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "update_iep_denied"); !ok {
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"title":      req.Title,
		"goal":       req.Goal,
		"plan":       req.Plan,
		"evaluation": req.Evaluation,
		"status":     req.Status,
		"update_by":  req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduCaseIEP{}).
		Where("id = ? and case_id = ?", c.Param("iepId"), c.Param("id")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("iepId"), "更新成功")
}

func (e EduCase) DeleteIEPs(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "delete_ieps_denied"); !ok {
		return
	}
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Delete(&models.EduCaseIEP{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCase) GetAssessments(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "view", "view_assessments_denied"); !ok {
		return
	}
	list := make([]models.EduCaseAssessment, 0)
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if shouldDesensitize(c) {
		desensitizeAssessments(list)
	}
	e.writeAccessLog(c, parsePathId(c.Param("id")), "view_assessments")
	e.OK(list, "查询成功")
}

func (e EduCase) InsertAssessment(c *gin.Context) {
	req := models.EduCaseAssessment{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "add_assessment_denied"); !ok {
		return
	}
	req.CaseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCase) UpdateAssessment(c *gin.Context) {
	req := models.EduCaseAssessment{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "update_assessment_denied"); !ok {
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"tool_name":   req.ToolName,
		"result":      req.Result,
		"assessed_at": req.AssessedAt,
		"update_by":   req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduCaseAssessment{}).
		Where("id = ? and case_id = ?", c.Param("assessmentId"), c.Param("id")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("assessmentId"), "更新成功")
}

func (e EduCase) DeleteAssessments(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "delete_assessments_denied"); !ok {
		return
	}
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Delete(&models.EduCaseAssessment{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCase) GetInterventions(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "view", "view_interventions_denied"); !ok {
		return
	}
	list := make([]models.EduCaseIntervention, 0)
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Order("id desc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if shouldDesensitize(c) {
		desensitizeInterventions(list)
	}
	e.writeAccessLog(c, parsePathId(c.Param("id")), "view_interventions")
	e.OK(list, "查询成功")
}

func (e EduCase) InsertIntervention(c *gin.Context) {
	req := models.EduCaseIntervention{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "add_intervention_denied"); !ok {
		return
	}
	req.CaseId = parsePathId(c.Param("id"))
	req.SetCreateBy(user.GetUserId(c))
	if req.Status == "" {
		req.Status = "active"
	}
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.Id, "创建成功")
}

func (e EduCase) UpdateIntervention(c *gin.Context) {
	req := models.EduCaseIntervention{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "update_intervention_denied"); !ok {
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	updates := map[string]interface{}{
		"title":      req.Title,
		"content":    req.Content,
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"status":     req.Status,
		"update_by":  req.UpdateBy,
	}
	if err := e.Orm.Model(&models.EduCaseIntervention{}).
		Where("id = ? and case_id = ?", c.Param("interventionId"), c.Param("id")).
		Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("interventionId"), "更新成功")
}

func (e EduCase) DeleteInterventions(c *gin.Context) {
	req := struct {
		Ids []int `json:"ids"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if _, ok := e.ensureCaseAccess(c, parsePathId(c.Param("id")), "edit", "delete_interventions_denied"); !ok {
		return
	}
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Delete(&models.EduCaseIntervention{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
