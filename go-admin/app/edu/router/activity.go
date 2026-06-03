package router

import (
	"go-admin/app/edu/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerActivityRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerPublicActivityRouter)
}

func registerPublicActivityRouter(v1 *gin.RouterGroup) {
	api := apis.EduActivity{}
	r := v1.Group("/portal/activities")
	{
		r.GET("", api.PublicGetPage)
		r.GET("/:id", api.PublicGet)
		r.GET("/:id/signup-state", api.PublicSignupState)
		r.POST("/:id/signup", api.PublicSignup)
		r.DELETE("/:id/signup", api.PublicCancelSignup)
		r.POST("/:id/checkin", api.PublicCheckin)
		r.POST("/:id/outcomes/files/upload", api.PublicUploadOutcomeFile)
		r.POST("/:id/outcomes", api.PublicSubmitOutcome)
	}
}

func registerActivityRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.EduActivity{}
	r := v1.Group("/edu/activities").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
		r.POST("/:id/signup", api.Signup)
		r.GET("/:id/signups", api.GetSignups)
		r.POST("/:id/signups", api.InsertSignup)
		r.PUT("/:id/signups/:signupId", api.UpdateSignup)
		r.DELETE("/:id/signups", api.DeleteSignups)
		r.GET("/:id/checkins", api.GetCheckins)
		r.POST("/:id/checkins", api.InsertCheckin)
		r.DELETE("/:id/checkins", api.DeleteCheckins)
		r.GET("/:id/outcomes", api.GetOutcomes)
		r.POST("/:id/outcomes", api.InsertOutcome)
		r.PUT("/:id/outcomes/:outcomeId", api.UpdateOutcome)
		r.DELETE("/:id/outcomes", api.DeleteOutcomes)
	}
}
