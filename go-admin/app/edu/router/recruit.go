package router

import (
	"go-admin/app/edu/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, registerPublicRecruitRouter)
	routerCheckRole = append(routerCheckRole, registerRecruitRouter)
}

func registerPublicRecruitRouter(v1 *gin.RouterGroup) {
	api := apis.EduRecruit{}
	r := v1.Group("/recruit")
	{
		r.GET("/jobs", api.PublicJobs)
		r.GET("/jobs/:id", api.PublicJob)
		r.GET("/companies", api.PublicCompanies)
		r.GET("/companies/:id", api.PublicCompany)
		r.POST("/company-applications", api.SubmitCompanyApplication)
		r.POST("/job-applications", api.SubmitJobApplication)
	}
}

func registerRecruitRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.EduRecruit{}
	r := v1.Group("/admin/recruit").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("/stats", api.Stats)
		r.GET("/reviews/companies", api.ReviewCompanies)
		r.GET("/reviews/jobs", api.ReviewJobs)
		r.GET("/reviews/:id", api.ReviewDetail)
		r.POST("/reviews/:id/approve", api.ApproveReview)
		r.POST("/reviews/:id/reject", api.RejectReview)
		r.GET("/companies", api.AdminCompanies)
		r.POST("/companies", api.CreateCompany)
		r.GET("/companies/:id", api.AdminCompany)
		r.PUT("/companies/:id", api.UpdateCompany)
		r.POST("/companies/:id/enable", api.EnableCompany)
		r.POST("/companies/:id/disable", api.DisableCompany)
		r.GET("/jobs", api.AdminJobs)
		r.POST("/jobs", api.CreateJob)
		r.GET("/jobs/:id", api.AdminJob)
		r.PUT("/jobs/:id", api.UpdateJob)
		r.POST("/jobs/:id/publish", api.PublishJob)
		r.POST("/jobs/:id/offline", api.OfflineJob)
		r.DELETE("/jobs/:id", api.DeleteJob)
	}
}
