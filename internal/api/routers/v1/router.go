package router_v1

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	RegisterCommonRouter(v1.Group(""))
	RegisterSysUserRouter(v1.Group("/user"))
	RegisterSysRoleRouter(v1.Group("/role"))
	RegisterSysMenuRouter(v1.Group("/menu"))
	RegisterSysCasbinRouter(v1.Group("/casbin"), router)
	RegisterPurRouter(v1.Group("/pur"))
	RegisterPursRouter(v1.Group("/purs"))
	RegisterPurdRouter(v1.Group("/purd"))
	RegisterPurSizeRouter(v1.Group("/purSize"))
	RegisterPurPRouter(v1.Group("/purp"))
	RegisterPurCFMRouter(v1.Group("/cfm"))
}
