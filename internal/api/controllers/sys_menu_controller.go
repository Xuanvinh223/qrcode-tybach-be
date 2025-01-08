package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jinzhu/copier"
)

type SysMenuController struct {
	*BaseController
}

var SysMenu = &SysMenuController{}

func (c *SysMenuController) Create(ctx *gin.Context) {
	var requestParams request.SysMenuCreateRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var menu entities.SysMenu
	if err := copier.Copy(&menu, requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	_ = services.SysMenu.Create(&menu)
	if err := services.SysMenu.Save(&menu); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *SysMenuController) View(ctx *gin.Context) {
	id := ctx.Param("id")
	var menu *entities.SysMenu
	if notFound, _ := services.SysMenu.FirstById(&menu, gconv.Uint64(id)); notFound {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "can not found the menu")
		return
	}

	response.OkWithData(ctx, menu)
}

func (c *SysMenuController) List(ctx *gin.Context) {
	var requestParams request.PageInfo
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var dataList []entities.SysMenu
	page, err := services.SysMenu.Pagination(&entities.SysMenu{}, &dataList, requestParams, []string{})
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, page)
}

func (c *SysMenuController) Tree(ctx *gin.Context) {
	menuTreeList, err := services.SysMenu.MenuTreeList()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, menuTreeList)
}

func (c *SysMenuController) Update(ctx *gin.Context) {
	var requestParams request.SysMenuUpdateRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var menu *entities.SysMenu
	if notFound, _ := services.SysMenu.FirstById(&menu, requestParams.ID); notFound {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "can not found the menu")
		return
	}

	if err := copier.Copy(&menu, requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	if err := services.SysMenu.Save(&menu); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, menu)
}

func (c *SysMenuController) Delete(ctx *gin.Context) {
	var requestParams request.SysMenuDeleteRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.SysMenu.Delete(requestParams.ID); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}
