package services

import (
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jinzhu/copier"
)

type SysMenuService struct {
	*BaseService
}

var SysMenu = &SysMenuService{}

func (s *SysMenuService) MenuTreeList() ([]types.SysMenuTreeList, error) {
	var menuList []entities.SysMenu
	var treeMenuList []types.SysMenuTreeList

	_, err := s.Scan(&entities.SysMenu{}, &entities.SysMenu{}, &menuList)
	if err != nil {
		return nil, err
	}
	if err := copier.Copy(&treeMenuList, menuList); err != nil {
		return nil, err
	}

	result := generateMenuTreeList(treeMenuList, "0")

	return result, nil
}

func (s *SysMenuService) Delete(ids []uint64) error {
	var err error

	for _, menuId := range ids {
		_, err = s.DeleteByID(&entities.SysMenu{}, menuId)
		_, err = s.DeleteByWhere(&entities.SysRoleMenu{}, entities.SysRoleMenu{MenuID: menuId})
	}

	if err != nil {
		return err
	}

	return nil
}

func generateMenuTreeList(menus []types.SysMenuTreeList, pid string) []types.SysMenuTreeList {
	var result []types.SysMenuTreeList

	for _, menu := range menus {
		if menu.MenuPid == pid {
			children := generateMenuTreeList(menus, gconv.String(menu.ID))
			if len(children) > 0 {
				menu.Children = children
			}
			result = append(result, menu)
		}
	}

	return result
}
