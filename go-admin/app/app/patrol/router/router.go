package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/core/global"
	"go-admin/core/runtime"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
	routerCheckRole   = make([]func(v1 *gin.RouterGroup), 0)
)

// InitRouter 初始化巡检模块路由
func InitRouter() {
	var r *gin.Engine
	h := runtime.RuntimeConfig.GetEngine()
	if h == nil {
		panic("not found engine...")
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		panic("not support other engine")
	}

	// 无需认证的路由(Unity 客户端上报)
	noCheckRoleRouter(r)
	// 需要认证的路由(React 管理端)
	checkRoleRouter(r)
}

func noCheckRoleRouter(r *gin.Engine) {
	v1 := r.Group(global.RouteRootPath + "/v1")
	for _, f := range routerNoCheckRole {
		f(v1)
	}
}

func checkRoleRouter(r *gin.Engine) {
	v1 := r.Group(global.RouteRootPath + "/v1")
	for _, f := range routerCheckRole {
		f(v1)
	}
}
