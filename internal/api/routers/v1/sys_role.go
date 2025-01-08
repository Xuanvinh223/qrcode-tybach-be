package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"
	"tyxuan-web-printlabel-api/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterSysRoleRouter(router *gin.RouterGroup) {
	router.POST("/create", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysRole.Create)
	router.GET("/view/:id", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysRole.View)
	router.GET("/list", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysRole.List)
	router.POST("/update", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysRole.Update)
	router.DELETE("/delete", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysRole.Delete)
}
