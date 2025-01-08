package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
)

type SysCasbinController struct {
	*BaseController
}

var SysCasbin = &SysCasbinController{}

func (c *SysCasbinController) Create(ctx *gin.Context) {
	var requestParams request.SysCasbinCreateRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	route := entities.SysCasbin{
		Ptype: "p",
		V0:    gconv.String(requestParams.RoleId),
		V1:    requestParams.Url,
		V2:    requestParams.Method,
	}
	_ = services.SysCasbin.Create(&route)
	if err := services.SysCasbin.Save(&route); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *SysCasbinController) View(ctx *gin.Context) {
	id := ctx.Param("id")
	var route *entities.SysCasbin
	if notFound, _ := services.SysCasbin.FirstById(&route, gconv.Uint64(id)); notFound {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "can not found the route")
		return
	}

	response.OkWithData(ctx, route)
}

func (c *SysCasbinController) List(ctx *gin.Context) {
	var requestParams request.PageInfo
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var dataList []entities.SysCasbin
	page, err := services.SysCasbin.Pagination(&entities.SysCasbin{}, &dataList, requestParams, []string{})
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, page)
}

func (c *SysCasbinController) Update(ctx *gin.Context) {
	var requestParams request.SysCasbinUpdateRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var route *entities.SysCasbin
	if notFound, _ := services.SysCasbin.FirstById(&route, requestParams.ID); notFound {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "can not found the route")
		return
	}

	route.V0 = gconv.String(requestParams.RoleId)
	route.V1 = requestParams.Url
	route.V2 = requestParams.Method
	if err := services.SysCasbin.Save(&route); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, route)
}

func (c *SysCasbinController) Delete(ctx *gin.Context) {
	var requestParams request.SysCasbinDeleteRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if _, err := services.SysCasbin.DeleteByIDS(&entities.SysCasbin{}, requestParams.ID); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *SysCasbinController) Routes(ctx *gin.Context) {
	var requestParams request.RouteRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	routerInstance := ctx.MustGet("RouterInstance").(*gin.Engine)
	casbinRuleList, err := services.SysCasbin.AddRoutes(requestParams, routerInstance)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	response.OkWithData(ctx, casbinRuleList)
}
