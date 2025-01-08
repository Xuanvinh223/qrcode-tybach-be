package entities

import (
	"time"

	"gorm.io/gorm"
)

type SysUser struct {
	BaseModel
	UserName  string    `gorm:"column:user_name;uniqueIndex;not null" json:"userName"`
	RealName  string    `gorm:"column:real_name;not null" json:"realName"`
	Password  string    `gorm:"column:password;not null" json:"password"`
	Pass      string    `gorm:"column:pass;not null" json:"pass"`
	Email     string    `gorm:"column:email;uniqueIndex;not null" json:"email"`
	LoginTime time.Time `gorm:"column:login_time;type:datetime" json:"loginTime"`
	LoginIp   string    `gorm:"column:login_ip" json:"loginIp"`
	State     int       `gorm:"column:state;default:1;not null" json:"state"` // 1:enable 2:disable
}

func (m *SysUser) BeforeCreate(*gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *SysUser) BeforeUpdate(*gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
