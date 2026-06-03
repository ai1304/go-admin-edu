package router

import (
	"go-admin/app/edu/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerNewsRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerPublicNewsRouter)
}

func registerPublicNewsRouter(v1 *gin.RouterGroup) {
	newsApi := apis.EduNews{}
	r := v1.Group("/portal/news")
	{
		r.GET("", newsApi.PublicGetPage)
		r.GET("/:id", newsApi.PublicGet)
		r.POST("/:id/like", newsApi.PublicLike)
	}
}

func registerNewsRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	newsApi := apis.EduNews{}
	r := v1.Group("/edu/news").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", newsApi.GetPage)
		r.GET("/:id", newsApi.Get)
		r.POST("", newsApi.Insert)
		r.PUT("/:id", newsApi.Update)
		r.DELETE("", newsApi.Delete)
	}
}
