package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPurCFMRouter(router *gin.RouterGroup) {
	router.GET("/PurCFM", controllers.PurCFM.GetPurCFMData)
	router.GET("/DetailCFM", controllers.PurCFM.GetDetailCFMData)
	router.POST("/CFMALL", controllers.PurCFM.CFMALL)
}
