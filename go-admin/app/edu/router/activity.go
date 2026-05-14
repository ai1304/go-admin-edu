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
	}
}
