package router

import (
	"go-admin/app/edu/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerExpertRouter, registerStatsRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerPublicExpertRouter)
}

func registerPublicExpertRouter(v1 *gin.RouterGroup) {
	api := apis.EduExpert{}
	r := v1.Group("/portal/experts")
	{
		r.GET("", api.PublicGetPage)
		r.GET("/:id", api.PublicGet)
		r.GET("/:id/resources/:resourceId/access-url", api.PublicResourceAccessURL)
		r.GET("/:id/favorite-state", api.PublicFavoriteState)
		r.POST("/:id/favorite", api.PublicFavorite)
		r.DELETE("/:id/favorite", api.PublicUnfavorite)
		r.PUT("/:id/share", api.PublicShare)
	}
}

func registerExpertRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.EduExpert{}
	r := v1.Group("/edu/experts").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
		r.GET("/:id/resources", api.GetResources)
		r.POST("/:id/resources", api.InsertResource)
		r.PUT("/:id/resources/:resourceId", api.UpdateResource)
		r.DELETE("/:id/resources", api.DeleteResources)
	}
}

func registerStatsRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.EduStats{}
	r := v1.Group("/edu/stats").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("/overview", api.Overview)
		r.GET("/resources", api.Resources)
		r.GET("/courses", api.Courses)
		r.GET("/activities", api.Activities)
		r.GET("/schools", api.Schools)
		r.GET("/teachers", api.Teachers)
		r.GET("/students", api.Students)
		r.GET("/cases", api.Cases)
		r.GET("/export", api.Export)
	}
}
