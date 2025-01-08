package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPurPRouter(router *gin.RouterGroup) {
	router.GET("/PurListPrint", controllers.PurP.GetPurListPrintData)
	router.POST("/DetailListPrint", controllers.PurP.GetDetailListPrintData)
	router.POST("/LabelPrint", controllers.PurP.GetLabelPrintData)
}
