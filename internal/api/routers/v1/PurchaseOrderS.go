package router_v1

import (
	"tyxuan-web-printlabel-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPursRouter(router *gin.RouterGroup) {
	router.GET("/PurListS", controllers.Purs.GetPurListDataS)
	router.GET("/DetailListS", controllers.Purs.GetDetailListDataS)
	router.POST("/SetPackQtyS", controllers.Purs.SetPackQtyS)
	router.POST("/AddPackS", controllers.Purs.AddPackS)
	router.POST("/DelPackS", controllers.Purs.DelPackS)
	router.POST("/LoadingQtyS", controllers.Purs.LoadingQtyS)
	router.POST("/CFMS", controllers.Purs.CFMS)
	router.POST("/UploadLotFileS", controllers.Purs.UploadLotFileS)
}
