package entities

import (
	"time"

	"gorm.io/gorm"
)

type Zszl_lot struct {
	ZSDH     string    `gorm:"column:ZSDH;primary_key;not null" json:"ZSDH"`
	LotNO    string    `gorm:"column:LotNO;primary_key;not null" json:"LotNO"`
	LotFile  string    `gorm:"column:LotFile;primary_key;not null" json:"LotFile"`
	USERDATE time.Time `gorm:"column:USERDATE" json:"USERDATE"`
}

func (Zszl_lot) TableName() string {
	return "zszl_lot"
}

func (m *Zszl_lot) BeforeCreate(*gorm.DB) error {
	m.USERDATE = time.Now()
	return nil
}

func (m *Zszl_lot) BeforeUpdate(*gorm.DB) error {
	m.USERDATE = time.Now()
	return nil
}
