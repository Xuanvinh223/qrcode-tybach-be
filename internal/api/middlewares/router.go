package middlewares

import (
	"github.com/gin-gonic/gin"
)

func PassRouterInstance(router *gin.Engine) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("RouterInstance", router)
		ctx.Next()
	}
}
