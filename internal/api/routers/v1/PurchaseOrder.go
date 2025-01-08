package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPurRouter(router *gin.RouterGroup) {
	router.GET("/PurList", controllers.Pur.GetPurListData)
	router.GET("/DetailList", controllers.Pur.GetDetailListData)
	router.POST("/SetPackQty", controllers.Pur.SetPackQty)
	router.POST("/ImportPackQty", controllers.Pur.Import_PackQty)
	router.POST("/AddPack", controllers.Pur.AddPack)
	router.POST("/DelPack", controllers.Pur.DelPack)
	router.POST("/LoadingQty", controllers.Pur.LoadingQty)
	router.POST("/CFM", controllers.Pur.CFM)
	router.POST("/UploadLotFile", controllers.Pur.UploadLotFile)
}
