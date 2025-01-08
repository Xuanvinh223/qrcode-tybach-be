package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type CommonController struct {
	*BaseController
}

var Common = &CommonController{}

func (c *CommonController) Ping(ctx *gin.Context) {
	response.Ok(ctx)
}

func (c *CommonController) Initialization(ctx *gin.Context) {
	if err := services.Common.Initialization(); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.Ok(ctx)
}
