package entities

import (
	"time"

	"gorm.io/gorm"
)

type SysRole struct {
	BaseModel
	RoleName    string `gorm:"column:role_name" json:"roleName"`
	Description string `gorm:"column:description" json:"description"`
}

func (m *SysRole) BeforeCreate(*gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *SysRole) BeforeUpdate(*gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
