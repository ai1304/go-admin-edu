package router

import (
	"go-admin/app/edu/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCaseRouter)
}

func registerCaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.EduCase{}
	r := v1.Group("/edu/cases").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
		r.GET("/:id/access-logs", api.GetAccessLogs)
		r.GET("/:id/authorizations", api.GetAuthorizations)
		r.POST("/:id/authorizations", api.AddAuthorization)
		r.PUT("/:id/authorizations/:authorizationId", api.UpdateAuthorization)
		r.DELETE("/:id/authorizations", api.DeleteAuthorizations)
		r.GET("/:id/ieps", api.GetIEPs)
		r.POST("/:id/ieps", api.AddIEP)
		r.PUT("/:id/ieps/:iepId", api.UpdateIEP)
		r.DELETE("/:id/ieps", api.DeleteIEPs)
		r.GET("/:id/assessments", api.GetAssessments)
		r.POST("/:id/assessments", api.InsertAssessment)
		r.PUT("/:id/assessments/:assessmentId", api.UpdateAssessment)
		r.DELETE("/:id/assessments", api.DeleteAssessments)
		r.GET("/:id/interventions", api.GetInterventions)
		r.POST("/:id/interventions", api.InsertIntervention)
		r.PUT("/:id/interventions/:interventionId", api.UpdateIntervention)
		r.DELETE("/:id/interventions", api.DeleteInterventions)
	}
}
