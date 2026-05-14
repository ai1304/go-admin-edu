package router

import (
	"go-admin/app/edu/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCourseRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerPublicCourseRouter)
}

func registerPublicCourseRouter(v1 *gin.RouterGroup) {
	api := apis.EduCourse{}
	r := v1.Group("/portal/courses")
	{
		r.GET("", api.PublicGetPage)
		r.GET("/:id", api.PublicGet)
	}
}

func registerCourseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.EduCourse{}
	r := v1.Group("/edu/courses").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
		r.GET("/:id/chapters", api.GetChapters)
		r.POST("/:id/chapters", api.InsertChapter)
		r.PUT("/:id/chapters/:chapterId", api.UpdateChapter)
		r.DELETE("/:id/chapters", api.DeleteChapters)
		r.GET("/:id/lessons", api.GetLessons)
		r.POST("/:id/lessons", api.InsertLesson)
		r.PUT("/:id/lessons/:lessonId", api.UpdateLesson)
		r.DELETE("/:id/lessons", api.DeleteLessons)
	}
}
