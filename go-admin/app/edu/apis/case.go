package apis

import (
	"go-admin/app/edu/models"
	"go-admin/common/dto"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

type EduCase struct {
	api.Api
}

type caseQuery struct {
	dto.Pagination
	Keyword     string `form:"keyword"`
	Status      string `form:"status"`
	SchoolId    int    `form:"schoolId"`
	Desensitize bool   `form:"desensitize"`
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

func (e EduCase) GetPage(c *gin.Context) {
	req := caseQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCase, 0)
	db := e.Orm.Model(&models.EduCase{})
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
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if err := db.Order("id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if req.Desensitize {
		desensitizeCases(list)
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e EduCase) Get(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduCase
	if err := e.Orm.First(&data, c.Param("id")).Error; err != nil {
		e.Error(500, err, "查询失败")
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
	if err := e.Orm.Delete(&models.EduCase{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}

func (e EduCase) GetAccessLogs(c *gin.Context) {
	req := caseAccessLogQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduCaseAccessLog, 0)
	db := e.Orm.Model(&models.EduCaseAccessLog{}).Where("case_id = ?", c.Param("id"))
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

func (e EduCase) GetAuthorizations(c *gin.Context) {
	req := caseAuthorizationQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
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
	if err := e.Orm.Where("case_id = ?", c.Param("id")).Delete(&models.EduCaseIntervention{}, req.Ids).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.Ids, "删除成功")
}
