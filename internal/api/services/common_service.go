package services

import (
	"fmt"
	"net/http"

	"tyxuan-web-printlabel-api/internal/pkg/config"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/pkg/crypto"
)

type CommonService struct {
	*BaseService
}

var Common = &CommonService{}

func (s *CommonService) Initialization() error {
	var err error

	// SysUser
	if notFound, _ := s.FirstById(&entities.SysUser{}, 1); notFound {
		user := entities.SysUser{
			UserName: "admin",
			RealName: "Administrator",
			Password: crypto.HashAndSalt([]byte("admin")),
			Email:    "admin@gmail.com",
			State:    1,
		}
		_ = s.Create(&user)
		err = s.Save(&user)
	}

	// SysRole
	if notFound, _ := s.FirstById(&entities.SysRole{}, 1); notFound {
		role := entities.SysRole{
			RoleName:    "Administrator",
			Description: "Root Administrator",
		}
		_ = s.Create(&role)
		err = s.Save(&role)
	}

	// SysUserRole
	if notFound, _ := s.FirstById(&entities.SysUserRole{}, 1); notFound {
		userRole := entities.SysUserRole{
			UserID: 1,
			RoleID: 1,
		}
		_ = s.Create(&userRole)
		err = s.Save(&userRole)
	}

	// SysMenu
	if notFound, _ := s.FirstById(&entities.SysMenu{}, 1); notFound {
		menus := []entities.SysMenu{
			{MenuName: "Home", MenuPid: "0", State: 1},
			{MenuName: "Permission", MenuPid: "0", State: 1},
			{MenuName: "PurchaseOrder", MenuPid: "0", State: 1},
			{MenuName: "QR_Code", MenuPid: "0", State: 1},
			{MenuName: "Users", MenuPid: "2", State: 1},
			{MenuName: "Roles", MenuPid: "2", State: 1},
			{MenuName: "Menus", MenuPid: "2", State: 1},
			{MenuName: "API", MenuPid: "2", State: 1},
			{MenuName: "PrintLabelQR", MenuPid: "2", State: 1},
		}
		for _, menu := range menus {
			_ = s.Create(&menu)
			err = s.Save(&menu)
		}
	}

	// SysRoleMenu
	if notFound, _ := s.FirstById(&entities.SysRoleMenu{}, 1); notFound {
		menuIds := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}
		for _, menuId := range menuIds {
			roleMenu := entities.SysRoleMenu{
				RoleID: 1,
				MenuID: menuId,
			}
			_ = s.Create(&roleMenu)
			err = s.Save(&roleMenu)
		}
	}

	// SysCasbin
	requestURL := fmt.Sprintf("http://localhost:%s/api/v1/casbin/routes?rtype=/api/v1&roleId=1", config.GetConfig().Server.Port)
	_, err = http.Get(requestURL)

	if err != nil {
		return err
	}

	return nil
}
