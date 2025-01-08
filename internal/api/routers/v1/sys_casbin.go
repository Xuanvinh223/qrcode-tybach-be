package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"
	"tyxuan-web-printlabel-api/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterSysCasbinRouter(router *gin.RouterGroup, routerInstance *gin.Engine) {
	router.GET("/routes", middlewares.PassRouterInstance(routerInstance), controllers.SysCasbin.Routes)
	router.POST("/create", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysCasbin.Create)
	router.GET("/view/:id", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysCasbin.View)
	router.GET("/list", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysCasbin.List)
	router.POST("/update", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysCasbin.Update)
	router.DELETE("/delete", middlewares.JWTAuth, middlewares.CasbinAuth, controllers.SysCasbin.Delete)
}
