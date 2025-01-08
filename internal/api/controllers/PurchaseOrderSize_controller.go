package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type PurSizeController struct {
	*BaseController
}

var PurSize = &PurSizeController{}

func (c *PurSizeController) GetPurListDataSize(ctx *gin.Context) {
	var requestParams request.PurListRequestSize
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.PurSize.GetPurListSize(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PurSizeController) SetPackQtySize(ctx *gin.Context) {
	var requestParams request.PackQtyRequestSize
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.PurSize.SetPackQtySize(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurSizeController) GetDetailListDataSize(ctx *gin.Context) {
	var requestParams request.DetailListRequestSize
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := services.PurSize.GetDetailListSize(requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *PurSizeController) AddPackSize(ctx *gin.Context) {
	var requestParams request.PackRequest_SizeS
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.PurSize.AddPackSize(requestParams.SCNO, requestParams.CLBH, requestParams.USERID, requestParams.QTYAD, requestParams.XXCC); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurSizeController) DelPackSize(ctx *gin.Context) {
	var requestParams request.PackRequestSize
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.PurSize.DelPackSize(requestParams.SCNO, requestParams.CLBH, requestParams.USERID, requestParams.QTYAD); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurSizeController) LoadingQtySize(ctx *gin.Context) {
	var requestParams request.LoadRequestSize
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.PurSize.LoadingQtySize(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurSizeController) UploadLotFileSize(ctx *gin.Context) {
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

	if err := services.PurSize.UploadLotFileSize(requestParams, LotFile); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *PurSizeController) CFMSize(ctx *gin.Context) {
	var requestParams request.CFMRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.PurSize.CFMSize(requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}
