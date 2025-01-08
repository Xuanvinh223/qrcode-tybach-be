package services

import (
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"
	"tyxuan-web-printlabel-api/pkg/jwt"
)

type SysUserService struct {
	*BaseService
}

var SysUser = &SysUserService{}

func (s *SysUserService) FindUserById(id uint64) (*entities.SysUser, error) {
	var user entities.SysUser
	_, err := s.FirstById(&user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SysUserService) FindUserByUsername(username string) (*entities.SysUser, error) {
	var user entities.SysUser
	where := entities.SysUser{UserName: username, State: 1}
	_, err := s.First(&where, &user, []string{})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SysUserService) GetUserRoleId(user *entities.SysUser) []uint64 {
	var (
		roles   []entities.SysUserRole
		roleIds []uint64
	)

	where := entities.SysUserRole{UserID: user.ID}
	if err := s.Find(&where, &roles, []string{}); err != nil {
		return roleIds
	}

	for _, role := range roles {
		roleIds = append(roleIds, role.RoleID)
	}

	return roleIds
}

func (s *SysUserService) GetUserInfo(userClaims *jwt.UserClaims) *types.UserInfo {
	var permissions []string
	for _, roleId := range userClaims.RoleID {
		if roleId == 0 {
			continue
		}

		var role entities.SysRole
		_, err := s.FirstById(&role, roleId)
		if err != nil {
			continue
		} else {
			permissions = append(permissions, role.RoleName)
		}
	}

	var user entities.SysUser
	_, _ = s.FirstById(&user, userClaims.UserID)

	userInfo := &types.UserInfo{
		Permissions: permissions,
		Username:    user.UserName,
		RealName:    user.RealName,
		Email:       user.Email,
		Avatar:      "https://upload.wikimedia.org/wikipedia/commons/thumb/5/59/User-avatar.svg/2048px-User-avatar.svg.png",
	}

	return userInfo
}

func (s *SysUserService) Delete(ids []uint64) error {
	var err error

	for _, userId := range ids {
		_, err = s.DeleteByID(&entities.SysUser{}, userId)
		_, err = s.DeleteByWhere(&entities.SysUserRole{}, entities.SysUserRole{UserID: userId})
	}

	if err != nil {
		return err
	}

	return nil
}
