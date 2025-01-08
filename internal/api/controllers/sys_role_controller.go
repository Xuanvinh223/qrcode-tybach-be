package controllers

import (
	"net/http"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jinzhu/copier"
)

type SysRoleController struct {
	*BaseController
}

var SysRole = &SysRoleController{}

func (c *SysRoleController) Create(ctx *gin.Context) {
	var requestParams request.SysRoleCreateRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var role entities.SysRole
	if err := copier.Copy(&role, requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	_ = services.SysRole.Create(&role)
	if err := services.SysRole.Save(&role); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	if err := services.SysRole.RebuildRoleUserAndRoleMenu(role.ID, requestParams.RoleUserList, requestParams.RoleMenuList); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *SysRoleController) View(ctx *gin.Context) {
	id := ctx.Param("id")
	var role *entities.SysRole
	if notFound, _ := services.SysRole.FirstById(&role, gconv.Uint64(id)); notFound {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "can not found the role")
		return
	}

	result, err := services.SysRole.GetRoleUserIdsAndRoleMenuIds(role.ID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, err.Error())
		return
	}

	var roleInfo types.SysRoleInfo
	if err := copier.Copy(&roleInfo, role); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	roleInfo.RoleUserList = result["users"]
	roleInfo.RoleMenuList = result["menus"]

	response.OkWithData(ctx, roleInfo)
}

func (c *SysRoleController) List(ctx *gin.Context) {
	var requestParams request.PageInfo
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var dataList []entities.SysRole
	page, err := services.SysRole.Pagination(&entities.SysRole{}, &dataList, requestParams, []string{})
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, page)
}

func (c *SysRoleController) Update(ctx *gin.Context) {
	var requestParams request.SysRoleUpdateRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var role *entities.SysRole
	if notFound, _ := services.SysRole.FirstById(&role, requestParams.ID); notFound {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "can not found the role")
		return
	}

	if err := copier.Copy(&role, requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	if err := services.SysRole.Save(&role); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	if err := services.SysRole.RebuildRoleUserAndRoleMenu(role.ID, requestParams.RoleUserList, requestParams.RoleMenuList); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, role)
}

func (c *SysRoleController) Delete(ctx *gin.Context) {
	var requestParams request.SysRoleDeleteRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.SysRole.Delete(requestParams.ID); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}
