package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCommonRouter(router *gin.RouterGroup) {
	router.GET("/ping", controllers.Common.Ping)
	router.POST("/initialization", controllers.Common.Initialization)
}
