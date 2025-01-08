package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
)

type PurCFMService struct {
	*BaseController
}

var PurCFM = &PurCFMService{}

func (c *PurCFMService) GetPurCFMData(ctx *gin.Context) {
	var requestParams request.PurCFMRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var result []types.PurCFM
	var err error

	result, err = services.PurCFM.GetPurCFM(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PurCFMService) GetDetailCFMData(ctx *gin.Context) {
	var requestParams request.DetailCFMRequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var result []types.DetailCFM
	var err error

	result, err = services.PurCFM.GetDetailCFM(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PurCFMService) CFMALL(ctx *gin.Context) {
	var requestParams request.CFMAllRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.PurCFM.CFMALL(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}
