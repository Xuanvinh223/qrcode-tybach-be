package services

import (
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"

	"github.com/gogf/gf/v2/util/gconv"
)

type SysRoleService struct {
	*BaseService
}

var SysRole = &SysRoleService{}

func (s *SysRoleService) RebuildRoleUserAndRoleMenu(roleId uint64, userIds []uint64, menuIds []uint64) error {
	if _, err := s.DeleteByWhere(&entities.SysUserRole{}, entities.SysUserRole{RoleID: roleId}); err != nil {
		return err
	}

	if _, err := s.DeleteByWhere(&entities.SysRoleMenu{}, entities.SysRoleMenu{RoleID: roleId}); err != nil {
		return err
	}

	for _, userId := range userIds {
		var userRole entities.SysUserRole
		userRole.RoleID = roleId
		userRole.UserID = userId
		_ = s.Create(&userRole)
		if err := s.Save(&userRole); err != nil {
			return err
		}
	}

	for _, menuId := range menuIds {
		var roleMenu entities.SysRoleMenu
		roleMenu.MenuID = menuId
		roleMenu.RoleID = roleId
		_ = s.Create(&roleMenu)
		if err := s.Save(&roleMenu); err != nil {
			return err
		}
	}

	return nil
}

func (s *SysRoleService) GetRoleUserIdsAndRoleMenuIds(roleId uint64) (map[string][]uint64, error) {
	var (
		userRoles []entities.SysUserRole
		roleMenus []entities.SysRoleMenu
		result    = make(map[string][]uint64)
	)

	if err := s.Find(entities.SysUserRole{RoleID: roleId}, &userRoles, []string{}); err != nil {
		return nil, err
	}
	for _, role := range userRoles {
		result["users"] = append(result["users"], role.UserID)
	}

	if err := s.Find(entities.SysRoleMenu{RoleID: roleId}, &roleMenus, []string{}); err != nil {
		return nil, err
	}
	for _, menu := range roleMenus {
		result["menus"] = append(result["menus"], menu.MenuID)
	}

	return result, nil
}

func (s *SysRoleService) Delete(ids []uint64) error {
	var err error

	for _, roleId := range ids {
		if roleId == 1 {
			continue
		}

		_, err = s.DeleteByID(&entities.SysRole{}, roleId)
		_, err = s.DeleteByWhere(&entities.SysUserRole{}, entities.SysUserRole{RoleID: roleId})
		_, err = s.DeleteByWhere(&entities.SysRoleMenu{}, entities.SysRoleMenu{RoleID: roleId})
		_, err = s.DeleteByWhere(&entities.SysCasbin{}, entities.SysCasbin{V0: gconv.String(roleId)})
	}

	if err != nil {
		return err
	}

	return nil
}
