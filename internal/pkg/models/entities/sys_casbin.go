package entities

import (
	"time"

	"gorm.io/gorm"
)

type SysCasbin struct {
	BaseModel
	Ptype string `gorm:"column:ptype;not null" json:"ptype"` // policy type
	V0    string `gorm:"column:v0;not null" json:"v0"`       // role id
	V1    string `gorm:"column:v1;not null" json:"v1"`       // api url
	V2    string `gorm:"column:v2;not null" json:"v2"`       // api method
	V3    string `gorm:"column:v3;not null" json:"v3"`
	V4    string `gorm:"column:v4;not null" json:"v4"`
	V5    string `gorm:"column:v5;not null" json:"v5"`
	V6    string `gorm:"column:v6;not null" json:"v6"`
	V7    string `gorm:"column:v7;not null" json:"v7"`
}

func (SysCasbin) TableName() string {
	return "casbin_rule"
}

func (m *SysCasbin) BeforeCreate(*gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *SysCasbin) BeforeUpdate(*gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
