package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"
	"tyxuan-web-printlabel-api/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterSysMenuRouter(router *gin.RouterGroup) {
	router.POST("/create", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysMenu.Create)
	router.GET("/view/:id", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysMenu.View)
	router.GET("/list", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysMenu.List)
	router.GET("/tree", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysMenu.Tree)
	router.POST("/update", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysMenu.Update)
	router.DELETE("/delete", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysMenu.Delete)
}
