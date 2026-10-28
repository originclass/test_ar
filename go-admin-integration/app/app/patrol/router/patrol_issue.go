package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/app/patrol/apis"
	"go-admin/core/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerPatrolIssueRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerNoCheckPatrolIssueRouter)
}

// registerPatrolIssueRouter 注册需要认证的路由(React 管理端使用)
func registerPatrolIssueRouter(v1 *gin.RouterGroup) {
	api := apis.PatrolIssue{}
	r := v1.Group("/app/patrol/issue").Use(middleware.Auth()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id/status", api.UpdateStatus)
		r.DELETE("", api.Delete)
	}
}

// registerNoCheckPatrolIssueRouter 注册免认证路由(Unity AR 客户端直接上报问题)。
// 注意:与认证组使用不同路径,避免 gin 同 method+path 重复注册 panic。
func registerNoCheckPatrolIssueRouter(v1 *gin.RouterGroup) {
	api := apis.PatrolIssue{}
	r := v1.Group("/app/patrol")
	{
		r.POST("/issue/report", api.Insert)
	}
}
