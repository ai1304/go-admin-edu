package router

import (
	"go-admin/app/edu/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerResourceRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerPublicResourceRouter)
}

func registerPublicResourceRouter(v1 *gin.RouterGroup) {
	resourceApi := apis.EduResource{}
	r := v1.Group("/portal/resources")
	{
		r.GET("", resourceApi.PublicGetPage)
		r.GET("/:id", resourceApi.PublicGet)
	}
}

func registerResourceRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	categoryApi := apis.EduResourceCategory{}
	tagApi := apis.EduResourceTag{}
	resourceApi := apis.EduResource{}
	fileApi := apis.EduResourceFile{}

	categories := v1.Group("/edu/resource-categories").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		categories.GET("", categoryApi.GetPage)
		categories.POST("", categoryApi.Insert)
		categories.PUT("/:id", categoryApi.Update)
		categories.DELETE("", categoryApi.Delete)
	}

	tags := v1.Group("/edu/resource-tags").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		tags.GET("", tagApi.GetPage)
		tags.POST("", tagApi.Insert)
		tags.PUT("/:id", tagApi.Update)
		tags.DELETE("", tagApi.Delete)
	}

	resources := v1.Group("/edu/resources").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		resources.GET("", resourceApi.GetPage)
		resources.GET("/:id", resourceApi.Get)
		resources.POST("", resourceApi.Insert)
		resources.PUT("/:id", resourceApi.Update)
		resources.DELETE("", resourceApi.Delete)
		resources.PUT("/:id/submit-review", resourceApi.SubmitReview)
		resources.PUT("/:id/review", resourceApi.Review)
	}

	files := v1.Group("/edu/resource-files").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		files.GET("", fileApi.GetPage)
		files.POST("", fileApi.Insert)
		files.DELETE("", fileApi.Delete)
	}
}
