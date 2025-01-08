package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"
	"tyxuan-web-printlabel-api/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterSysUserRouter(router *gin.RouterGroup) {
	// backend
	router.POST("/create", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysUser.Create)
	router.GET("/view/:id", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysUser.View)
	router.GET("/list", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysUser.List)
	router.POST("/update", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysUser.Update)
	router.DELETE("/delete", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysUser.Delete)

	// common
	router.POST("/login_by_password", controllers.SysUser.LoginByPassword)
	router.GET("/info", middlewares.JWTAuth, controllers.SysUser.UserInfo)
}
