package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPurSizeRouter(router *gin.RouterGroup) {
	router.GET("/PurListSize", controllers.PurSize.GetPurListDataSize)
	router.GET("/DetailListSize", controllers.PurSize.GetDetailListDataSize)
	router.POST("/SetPackQtySize", controllers.PurSize.SetPackQtySize)
	router.POST("/AddPackSize", controllers.PurSize.AddPackSize)
	router.POST("/DelPackSize", controllers.PurSize.DelPackSize)
	router.POST("/LoadingQtySize", controllers.PurSize.LoadingQtySize)
	router.POST("/CFMSize", controllers.PurSize.CFMSize)
	router.POST("/UploadLotFileSize", controllers.PurSize.UploadLotFileSize)
}
