package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPurdRouter(router *gin.RouterGroup) {
	router.GET("/PurListD2", controllers.Purd.GetPurListDataD2)
	router.GET("/DetailListD2", controllers.Purd.GetDetailListDataD2)
	router.POST("/SetPackQtyD2", controllers.Purd.SetPackQtyD2)
	router.POST("/AddPackD2", controllers.Purd.AddPackD2)
	router.POST("/DelPackD2", controllers.Purd.DelPackD2)
	router.POST("/LoadingQtyD2", controllers.Purd.LoadingQtyD2)
	router.POST("/CFMD2", controllers.Purd.CFMD2)
	router.POST("/UploadLotFileD2", controllers.Purd.UploadLotFileD2)
}
