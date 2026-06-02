package apis

import (
	"go-admin/app/edu/models"
	"go-admin/common/dto"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"gorm.io/gorm"
)

type EduRecruit struct {
	api.Api
}

type recruitCompanyQuery struct {
	dto.Pagination
	Keyword       string `form:"keyword"`
	CompanyNature string `form:"companyNature"`
	Industry      string `form:"industry"`
	CompanySize   string `form:"companySize"`
	Region        string `form:"region"`
	Status        string `form:"status"`
	HasJobs       string `form:"hasJobs"`
}

type recruitJobQuery struct {
	dto.Pagination
	Keyword     string `form:"keyword"`
	JobType     string `form:"jobType"`
	Industry    string `form:"industry"`
	Location    string `form:"location"`
	Education   string `form:"education"`
	SalaryRange string `form:"salaryRange"`
	Status      string `form:"status"`
	Sort        string `form:"sort"`
	CompanyId   int    `form:"companyId"`
}

func (e EduRecruit) PublicJobs(c *gin.Context) {
	req := recruitJobQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	expireRecruitJobs(e.Orm)
	list := make([]models.EduRecruitJob, 0)
	now := time.Now().Format("2006-01-02")
	db := e.Orm.Model(&models.EduRecruitJob{}).Where("status = ? and review_status = ? and (deadline = '' or deadline >= ?)", "published", "approved", now)
	db = applyRecruitJobFilters(db, req)
	if err := pageRecruitJobs(db, req, &list, e); err != nil {
		return
	}
}

func (e EduRecruit) PublicJob(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	now := time.Now().Format("2006-01-02")
	var job models.EduRecruitJob
	if err := e.Orm.Where("status = ? and review_status = ? and (deadline = '' or deadline >= ?)", "published", "approved", now).First(&job, c.Param("id")).Error; err != nil {
		e.Error(404, err, "岗位不存在或已过期")
		return
	}
	var company models.EduRecruitCompany
	_ = e.Orm.Where("id = ? and status = ?", job.CompanyId, "normal").First(&company).Error
	e.OK(gin.H{"job": job, "company": company}, "查询成功")
}

func (e EduRecruit) PublicCompanies(c *gin.Context) {
	req := recruitCompanyQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduRecruitCompany, 0)
	db := e.Orm.Model(&models.EduRecruitCompany{}).Where("status = ? and review_status = ?", "normal", "approved")
	db = applyRecruitCompanyFilters(db, req)
	if req.HasJobs == "1" {
		db = db.Where("id in (select company_id from edu_recruit_job where deleted_at is null and status = ? and review_status = ?)", "published", "approved")
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
	fillRecruitCompanyJobCounts(e.Orm, list)
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e EduRecruit) PublicCompany(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var company models.EduRecruitCompany
	if err := e.Orm.Where("status = ? and review_status = ?", "normal", "approved").First(&company, c.Param("id")).Error; err != nil {
		e.Error(404, err, "企业不存在")
		return
	}
	jobs := make([]models.EduRecruitJob, 0)
	now := time.Now().Format("2006-01-02")
	_ = e.Orm.Where("company_id = ? and status = ? and review_status = ? and (deadline = '' or deadline >= ?)", company.Id, "published", "approved", now).Order("publish_time desc,id desc").Find(&jobs).Error
	e.OK(gin.H{"company": company, "jobs": jobs}, "查询成功")
}

func (e EduRecruit) SubmitCompanyApplication(c *gin.Context) {
	req := models.EduRecruitCompany{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.Status = "pending"
	req.ReviewStatus = "pending"
	req.CertStatus = "pending"
	req.SetCreateBy(user.GetUserId(c))
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "提交失败")
		return
	}
	e.OK(gin.H{"id": req.Id, "status": req.ReviewStatus}, "入驻申请已提交，请等待平台审核。")
}

func (e EduRecruit) SubmitJobApplication(c *gin.Context) {
	req := models.EduRecruitJob{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var company models.EduRecruitCompany
	if err := e.Orm.Where("id = ? and status = ? and review_status = ?", req.CompanyId, "normal", "approved").First(&company).Error; err != nil {
		e.Error(400, err, "企业入驻审核通过后，可提交岗位发布信息。")
		return
	}
	req.CompanyName = company.CompanyName
	req.Industry = company.Industry
	req.Status = "pending"
	req.ReviewStatus = "pending"
	req.SetCreateBy(user.GetUserId(c))
	if err := e.Orm.Create(&req).Error; err != nil {
		e.Error(500, err, "提交失败")
		return
	}
	e.OK(gin.H{"id": req.Id, "status": req.ReviewStatus}, "岗位信息已提交，请等待平台审核。")
}

func (e EduRecruit) AdminCompanies(c *gin.Context) {
	req := recruitCompanyQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduRecruitCompany, 0)
	db := applyEduUserScope(c, e.Orm.Model(&models.EduRecruitCompany{}))
	db = applyRecruitCompanyFilters(db, req)
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
	fillRecruitCompanyJobCounts(e.Orm, list)
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e EduRecruit) AdminCompany(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var company models.EduRecruitCompany
	if err := e.Orm.First(&company, c.Param("id")).Error; err != nil {
		e.Error(404, err, "企业不存在")
		return
	}
	jobs := make([]models.EduRecruitJob, 0)
	reviews := make([]models.EduRecruitReview, 0)
	_ = e.Orm.Where("company_id = ?", company.Id).Order("id desc").Find(&jobs).Error
	_ = e.Orm.Where("target_type = ? and target_id = ?", "company", company.Id).Order("id desc").Find(&reviews).Error
	e.OK(gin.H{"company": company, "jobs": jobs, "reviews": reviews}, "查询成功")
}

func (e EduRecruit) UpdateCompany(c *gin.Context) {
	req := models.EduRecruitCompany{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduRecruitCompany{}).Where("id = ?", c.Param("id")).Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("id"), "更新成功")
}

func (e EduRecruit) EnableCompany(c *gin.Context) {
	e.setCompanyStatus(c, "normal")
}

func (e EduRecruit) DisableCompany(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	id := parsePathId(c.Param("id"))
	err := e.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.EduRecruitCompany{}).Where("id = ?", id).Updates(map[string]any{"status": "disabled", "update_by": user.GetUserId(c)}).Error; err != nil {
			return err
		}
		return tx.Model(&models.EduRecruitJob{}).Where("company_id = ?", id).Updates(map[string]any{"status": "offline", "update_by": user.GetUserId(c)}).Error
	})
	if err != nil {
		e.Error(500, err, "禁用失败")
		return
	}
	e.OK(id, "禁用成功")
}

func (e EduRecruit) AdminJobs(c *gin.Context) {
	req := recruitJobQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	expireRecruitJobs(e.Orm)
	list := make([]models.EduRecruitJob, 0)
	db := applyEduUserScope(c, e.Orm.Model(&models.EduRecruitJob{}))
	db = applyRecruitJobFilters(db, req)
	if req.Status != "" {
		if req.Status == "expired" {
			db = db.Where("status = ? or (deadline <> '' and deadline < ?)", "expired", time.Now().Format("2006-01-02"))
		} else {
			db = db.Where("status = ?", req.Status)
		}
	}
	if err := pageRecruitJobs(db, req, &list, e); err != nil {
		return
	}
}

func (e EduRecruit) AdminJob(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var job models.EduRecruitJob
	if err := e.Orm.First(&job, c.Param("id")).Error; err != nil {
		e.Error(404, err, "岗位不存在")
		return
	}
	var company models.EduRecruitCompany
	reviews := make([]models.EduRecruitReview, 0)
	_ = e.Orm.First(&company, job.CompanyId).Error
	_ = e.Orm.Where("target_type = ? and target_id = ?", "job", job.Id).Order("id desc").Find(&reviews).Error
	e.OK(gin.H{"job": job, "company": company, "reviews": reviews}, "查询成功")
}

func (e EduRecruit) UpdateJob(c *gin.Context) {
	req := models.EduRecruitJob{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	if err := e.Orm.Model(&models.EduRecruitJob{}).Where("id = ?", c.Param("id")).Updates(&req).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(c.Param("id"), "更新成功")
}

func (e EduRecruit) PublishJob(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if err := e.Orm.Model(&models.EduRecruitJob{}).Where("id = ?", c.Param("id")).Updates(map[string]any{"status": "published", "publish_time": now, "update_by": user.GetUserId(c)}).Error; err != nil {
		e.Error(500, err, "上架失败")
		return
	}
	e.OK(c.Param("id"), "上架成功")
}

func (e EduRecruit) OfflineJob(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Model(&models.EduRecruitJob{}).Where("id = ?", c.Param("id")).Updates(map[string]any{"status": "offline", "update_by": user.GetUserId(c)}).Error; err != nil {
		e.Error(500, err, "下架失败")
		return
	}
	e.OK(c.Param("id"), "下架成功")
}

func (e EduRecruit) DeleteJob(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Delete(&models.EduRecruitJob{}, parsePathId(c.Param("id"))).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(c.Param("id"), "删除成功")
}

func (e EduRecruit) ReviewCompanies(c *gin.Context) {
	e.reviewList(c, "company")
}

func (e EduRecruit) ReviewJobs(c *gin.Context) {
	e.reviewList(c, "job")
}

func (e EduRecruit) ReviewDetail(c *gin.Context) {
	targetType := c.Query("type")
	id := parsePathId(c.Param("id"))
	if targetType == "job" {
		c.Params = append(c.Params, gin.Param{Key: "id", Value: c.Param("id")})
		e.AdminJob(c)
		return
	}
	if id == 0 {
		e.Error(400, nil, "缺少审核对象")
		return
	}
	e.AdminCompany(c)
}

func (e EduRecruit) ApproveReview(c *gin.Context) {
	e.handleReview(c, "approved")
}

func (e EduRecruit) RejectReview(c *gin.Context) {
	e.handleReview(c, "rejected")
}

func (e EduRecruit) Stats(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	expireRecruitJobs(e.Orm)
	result := map[string]int64{}
	result["pendingCompanies"] = countRecruit(e.Orm, &models.EduRecruitCompany{}, "review_status = ?", "pending")
	result["pendingJobs"] = countRecruit(e.Orm, &models.EduRecruitJob{}, "review_status = ?", "pending")
	result["companies"] = countRecruit(e.Orm, &models.EduRecruitCompany{}, "status <> ?", "")
	result["normalCompanies"] = countRecruit(e.Orm, &models.EduRecruitCompany{}, "status = ?", "normal")
	result["disabledCompanies"] = countRecruit(e.Orm, &models.EduRecruitCompany{}, "status = ?", "disabled")
	result["jobs"] = countRecruit(e.Orm, &models.EduRecruitJob{}, "status <> ?", "")
	result["publishedJobs"] = countRecruit(e.Orm, &models.EduRecruitJob{}, "status = ?", "published")
	result["offlineJobs"] = countRecruit(e.Orm, &models.EduRecruitJob{}, "status = ?", "offline")
	result["expiredJobs"] = countRecruit(e.Orm, &models.EduRecruitJob{}, "deadline <> '' and deadline < ?", time.Now().Format("2006-01-02"))
	today := time.Now().Format("2006-01-02")
	result["todayApplications"] = countRecruit(e.Orm, &models.EduRecruitCompany{}, "created_at >= ?", today) + countRecruit(e.Orm, &models.EduRecruitJob{}, "created_at >= ?", today)
	e.OK(result, "查询成功")
}

func (e EduRecruit) setCompanyStatus(c *gin.Context, status string) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if err := e.Orm.Model(&models.EduRecruitCompany{}).Where("id = ?", c.Param("id")).Updates(map[string]any{"status": status, "update_by": user.GetUserId(c)}).Error; err != nil {
		e.Error(500, err, "操作失败")
		return
	}
	e.OK(c.Param("id"), "操作成功")
}

func (e EduRecruit) reviewList(c *gin.Context, targetType string) {
	req := struct {
		dto.Pagination
		Keyword string `form:"keyword"`
		Status  string `form:"status"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	if targetType == "job" {
		list := make([]models.EduRecruitJob, 0)
		db := e.Orm.Model(&models.EduRecruitJob{})
		if req.Status != "" {
			db = db.Where("review_status = ?", req.Status)
		}
		if req.Keyword != "" {
			like := "%" + req.Keyword + "%"
			db = db.Where("job_name like ? or company_name like ? or contact_name like ?", like, like, like)
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
		return
	}
	list := make([]models.EduRecruitCompany, 0)
	db := e.Orm.Model(&models.EduRecruitCompany{})
	if req.Status != "" {
		db = db.Where("review_status = ?", req.Status)
	}
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("company_name like ? or credit_code like ? or contact_name like ?", like, like, like)
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

func (e EduRecruit) handleReview(c *gin.Context, status string) {
	req := struct {
		TargetType string `json:"targetType"`
		Opinion    string `json:"opinion"`
		Reason     string `json:"reason"`
	}{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.TargetType == "" {
		req.TargetType = c.Query("type")
	}
	if status == "rejected" && strings.TrimSpace(req.Opinion+req.Reason) == "" {
		e.Error(400, nil, "驳回必须填写原因")
		return
	}
	id := parsePathId(c.Param("id"))
	now := time.Now().Format("2006-01-02 15:04:05")
	err := e.Orm.Transaction(func(tx *gorm.DB) error {
		review := models.EduRecruitReview{TargetType: req.TargetType, TargetId: id, Action: status, Status: status, Opinion: req.Opinion, Reason: req.Reason}
		review.SetCreateBy(user.GetUserId(c))
		if err := tx.Create(&review).Error; err != nil {
			return err
		}
		updates := map[string]any{"review_status": status, "review_opinion": strings.TrimSpace(req.Reason + " " + req.Opinion), "reviewed_at": now, "update_by": user.GetUserId(c)}
		if req.TargetType == "job" {
			if status == "approved" {
				updates["status"] = "published"
				updates["publish_time"] = now
			} else {
				updates["status"] = "rejected"
			}
			return tx.Model(&models.EduRecruitJob{}).Where("id = ?", id).Updates(updates).Error
		}
		if status == "approved" {
			updates["status"] = "normal"
			updates["cert_status"] = "approved"
		} else {
			updates["status"] = "rejected"
			updates["cert_status"] = "rejected"
		}
		return tx.Model(&models.EduRecruitCompany{}).Where("id = ?", id).Updates(updates).Error
	})
	if err != nil {
		e.Error(500, err, "审核失败")
		return
	}
	e.OK(id, "审核成功")
}

func applyRecruitCompanyFilters(db *gorm.DB, req recruitCompanyQuery) *gorm.DB {
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("company_name like ? or contact_name like ? or industry like ?", like, like, like)
	}
	if req.CompanyNature != "" {
		db = db.Where("company_nature = ?", req.CompanyNature)
	}
	if req.Industry != "" {
		db = db.Where("industry = ?", req.Industry)
	}
	if req.CompanySize != "" {
		db = db.Where("company_size = ?", req.CompanySize)
	}
	if req.Region != "" {
		db = db.Where("region like ?", "%"+req.Region+"%")
	}
	return db
}

func applyRecruitJobFilters(db *gorm.DB, req recruitJobQuery) *gorm.DB {
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("job_name like ? or company_name like ? or major_direction like ?", like, like, like)
	}
	if req.JobType != "" {
		db = db.Where("job_type = ?", req.JobType)
	}
	if req.Industry != "" {
		db = db.Where("industry = ?", req.Industry)
	}
	if req.Location != "" {
		db = db.Where("location like ?", "%"+req.Location+"%")
	}
	if req.Education != "" {
		db = db.Where("education = ?", req.Education)
	}
	if req.SalaryRange != "" {
		db = db.Where("salary_range = ?", req.SalaryRange)
	}
	if req.CompanyId != 0 {
		db = db.Where("company_id = ?", req.CompanyId)
	}
	return db
}

func pageRecruitJobs(db *gorm.DB, req recruitJobQuery, list *[]models.EduRecruitJob, e EduRecruit) error {
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "查询失败")
		return err
	}
	order := "publish_time desc,id desc"
	if req.Sort == "deadline" {
		order = "deadline asc,id desc"
	} else if req.Sort == "headcount" {
		order = "headcount desc,id desc"
	}
	if err := db.Order(order).Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return err
	}
	e.PageOK(*list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	return nil
}

func fillRecruitCompanyJobCounts(db *gorm.DB, list []models.EduRecruitCompany) {
	for i := range list {
		_ = db.Model(&models.EduRecruitJob{}).Where("company_id = ? and status = ? and review_status = ?", list[i].Id, "published", "approved").Count(&list[i].JobCount).Error
	}
}

func countRecruit(db *gorm.DB, model any, query string, args ...any) int64 {
	var count int64
	_ = db.Model(model).Where(query, args...).Count(&count).Error
	return count
}

func expireRecruitJobs(db *gorm.DB) {
	today := time.Now().Format("2006-01-02")
	_ = db.Model(&models.EduRecruitJob{}).
		Where("status = ? and deadline <> '' and deadline < ?", "published", today).
		Updates(map[string]any{"status": "expired"}).Error
}
