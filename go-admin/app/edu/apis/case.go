package apis

import (
	"go-admin/app/edu/models"
	"go-admin/common/dto"

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
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
	SchoolId int    `form:"schoolId"`
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
