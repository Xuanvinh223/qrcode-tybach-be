package entities

import (
	"time"

	"gorm.io/gorm"
)

type SysRoleMenu struct {
	BaseModel
	RoleID uint64 `gorm:"column:role_id" json:"roleId"`
	MenuID uint64 `gorm:"column:menu_id" json:"menuId"`
}

func (m *SysRoleMenu) BeforeCreate(*gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *SysRoleMenu) BeforeUpdate(*gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
