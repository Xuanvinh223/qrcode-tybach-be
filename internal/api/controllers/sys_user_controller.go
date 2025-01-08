package controllers

import (
	"net/http"
	"time"

	"tyxuan-web-printlabel-api/internal/api/services"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"
	"tyxuan-web-printlabel-api/pkg/crypto"
	"tyxuan-web-printlabel-api/pkg/jwt"
	"tyxuan-web-printlabel-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	jwt_v4 "github.com/golang-jwt/jwt/v4"
)

type SysUserController struct {
	*BaseController
}

var SysUser = &SysUserController{}

func (c *SysUserController) Create(ctx *gin.Context) {
	var requestParams request.SysUserCreateRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	user := entities.SysUser{
		UserName:  requestParams.UserName,
		RealName:  requestParams.RealName,
		Password:  crypto.HashAndSalt([]byte(requestParams.Password)),
		Pass:      requestParams.Password,
		Email:     requestParams.Email,
		State:     requestParams.State,
		LoginTime: time.Now(),
	}

	_ = services.SysUser.Create(&user)
	if err := services.SysUser.Save(&user); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *SysUserController) View(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := services.SysUser.FindUserById(gconv.Uint64(id))
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "can not found the user")
		return
	}
	response.OkWithData(ctx, user)
}

func (c *SysUserController) List(ctx *gin.Context) {
	var requestParams request.PageInfo
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	var dataList []entities.SysUser
	page, err := services.SysUser.Pagination(&entities.SysUser{}, &dataList, requestParams, []string{})
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, page)
}

func (c *SysUserController) Update(ctx *gin.Context) {
	var requestParams request.SysUserUpdateRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	user, err := services.SysUser.FindUserById(requestParams.ID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, "can not found the user")
		return
	}

	user.RealName = requestParams.RealName
	user.Email = requestParams.Email
	user.State = requestParams.State
	if requestParams.Password != "" {
		user.Password = crypto.HashAndSalt([]byte(requestParams.Password))
		user.Pass = requestParams.Password
	}
	if err := services.SysUser.Save(&user); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, user)
}

func (c *SysUserController) Delete(ctx *gin.Context) {
	var requestParams request.SysUserDeleteRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := services.SysUser.Delete(requestParams.ID); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.Ok(ctx)
}

func (c *SysUserController) LoginByPassword(ctx *gin.Context) {
	var requestParams request.SysUserLoginByPasswordRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	user, err := services.SysUser.FindUserByUsername(requestParams.Username)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "can not found the user")
		return
	}

	isAuthenticated := crypto.ComparePassword(user.Password, []byte(requestParams.Password))
	if !isAuthenticated {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "password is incorrect")
		return
	}

	user.LoginTime = time.Now()
	user.LoginIp = ctx.ClientIP()
	if err := services.SysUser.Save(&user); err != nil {
		logger.Errorf("LoginByPassword: can not update user login time and login ip")
	}

	var tokenInfo types.TokenInfo
	tokenString := jwt.GenerateToken(&jwt.UserClaims{
		UserID:           user.ID,
		RoleID:           services.SysUser.GetUserRoleId(user),
		RegisteredClaims: jwt_v4.RegisteredClaims{},
	})
	tokenInfo.AccessToken = tokenString

	response.OkWithData(ctx, tokenInfo)
}

func (c *SysUserController) UserInfo(ctx *gin.Context) {
	userClaims := ctx.MustGet("authorization_payload").(*jwt.UserClaims)
	userInfo := services.SysUser.GetUserInfo(userClaims)

	response.OkWithData(ctx, userInfo)
}
