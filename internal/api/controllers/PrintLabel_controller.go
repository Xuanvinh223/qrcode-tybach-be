package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
)

type PurPrintService struct {
	*BaseController
}

var PurP = &PurPrintService{}

func (c *PurPrintService) GetPurListPrintData(ctx *gin.Context) {
	var requestParams request.PurListPrintRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var result []types.PurListPrint
	var err error

	result, err = services.PurP.GetPurListPrint(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PurPrintService) GetDetailListPrintData(ctx *gin.Context) {
	var requestParams request.DetailListPrintRequest1

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var result []types.DetailListPrint
	var err error

	result, err = services.PurP.GetDetailListPrint(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PurPrintService) GetLabelPrintData(ctx *gin.Context) {
	var requestParams request.LabelPrintRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var result []types.LabelPrint
	var err error

	result, err = services.PurP.GetLabelPrint(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}
