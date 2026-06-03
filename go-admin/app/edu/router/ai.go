package router

import (
	"go-admin/app/edu/apis"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerAIRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerPublicAIRouter)
}

func registerPublicAIRouter(v1 *gin.RouterGroup) {
	aiApi := apis.EduAI{}
	r := v1.Group("/portal/ai")
	{
		r.POST("/chat", aiApi.Chat)
		r.GET("/conversations", aiApi.MyConversations)
		r.GET("/conversations/:id", aiApi.ConversationDetail)
		r.DELETE("/conversations/:id", aiApi.PublicDeleteConversation)
	}
}

func registerAIRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	aiApi := apis.EduAI{}
	r := v1.Group("/edu/ai").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("/stats", aiApi.AdminStats)
		r.GET("/conversations", aiApi.AdminConversations)
		r.GET("/conversations/:id", aiApi.ConversationDetail)
		r.DELETE("/conversations/:id", aiApi.DeleteConversation)
	}
}
