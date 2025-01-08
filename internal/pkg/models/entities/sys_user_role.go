package entities

import (
	"time"

	"gorm.io/gorm"
)

type SysUserRole struct {
	BaseModel
	UserID uint64 `gorm:"column:user_id" json:"userId"`
	RoleID uint64 `gorm:"column:role_id" json:"roleId"`
}

func (m *SysUserRole) BeforeCreate(*gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *SysUserRole) BeforeUpdate(*gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
