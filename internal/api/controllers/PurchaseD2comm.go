package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type PurdController struct {
	*BaseController
}

var Purd = &PurdController{}

func (c *PurdController) GetPurListDataD2(ctx *gin.Context) {
	var requestParams request.PurListRequestKHPO
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.Purd.GetPurListD2(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PurdController) SetPackQtyD2(ctx *gin.Context) {
	var requestParams request.PackQtyRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purd.SetPackQtyD2(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurdController) GetDetailListDataD2(ctx *gin.Context) {
	var requestParams request.DetailListRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.Purd.GetDetailListD2(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PurdController) AddPackD2(ctx *gin.Context) {
	var requestParams request.PackRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purd.AddPackD2(requestParams.SCNO, requestParams.CLBH, requestParams.ZLBH, requestParams.USERID, requestParams.QTYAD); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurdController) DelPackD2(ctx *gin.Context) {
	var requestParams request.PackRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purd.DelPackD2(requestParams.SCNO, requestParams.CLBH, requestParams.USERID, requestParams.QTYAD); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurdController) LoadingQtyD2(ctx *gin.Context) {
	var requestParams request.LoadRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purd.LoadingQtyD2(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurdController) UploadLotFileD2(ctx *gin.Context) {
	var requestParams request.UploadLotFileRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	LotFile, err := ctx.FormFile("LotFile")
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	//ctx.SaveUploadedFile(LotFile, "./uploadfile/"+LotFile.Filename)

	if err := services.Purd.UploadLotFileD2(requestParams, LotFile); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurdController) CFMD2(ctx *gin.Context) {
	var requestParams request.CFMRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purd.CFMD2(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}
