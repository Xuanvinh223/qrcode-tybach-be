package entities

import (
	"time"

	"gorm.io/gorm"
)

type SysMenu struct {
	BaseModel
	MenuName    string `gorm:"column:menu_name" json:"menuName"`
	MenuPid     string `gorm:"column:menu_pid;not null" json:"menuPid"`
	Url         string `gorm:"column:url" json:"url"`
	Sort        string `gorm:"column:sort" json:"sort"`
	Description string `gorm:"column:description" json:"description"`
	State       int    `gorm:"column:state;default:1;not null" json:"state"` // 1:enable 2:disable
}

func (m *SysMenu) BeforeCreate(*gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *SysMenu) BeforeUpdate(*gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
