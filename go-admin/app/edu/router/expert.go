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
	}
}

func registerStatsRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.EduStats{}
	r := v1.Group("/edu/stats").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("/overview", api.Overview)
	}
}
