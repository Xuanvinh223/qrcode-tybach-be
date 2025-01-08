package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type PursController struct {
	*BaseController
}

var Purs = &PursController{}

func (c *PursController) GetPurListDataS(ctx *gin.Context) {
	var requestParams request.PurListRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.Purs.GetPurListS(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PursController) SetPackQtyS(ctx *gin.Context) {
	var requestParams request.PackQtyRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purs.SetPackQtyS(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PursController) GetDetailListDataS(ctx *gin.Context) {
	var requestParams request.DetailListRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.Purs.GetDetailListS(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PursController) AddPackS(ctx *gin.Context) {
	var requestParams request.PackRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purs.AddPackS(requestParams.SCNO, requestParams.CLBH, requestParams.ZLBH, requestParams.USERID, requestParams.QTYAD); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PursController) DelPackS(ctx *gin.Context) {
	var requestParams request.PackRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purs.DelPackS(requestParams.SCNO, requestParams.CLBH, requestParams.USERID, requestParams.QTYAD); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PursController) LoadingQtyS(ctx *gin.Context) {
	var requestParams request.LoadRequestS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purs.LoadingQtyS(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PursController) UploadLotFileS(ctx *gin.Context) {
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

	if err := services.Purs.UploadLotFileS(requestParams, LotFile); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PursController) CFMS(ctx *gin.Context) {
	var requestParams request.CFMRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.Purs.CFMS(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}
