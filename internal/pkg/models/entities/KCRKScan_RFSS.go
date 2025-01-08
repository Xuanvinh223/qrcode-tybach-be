package entities

import (
	"time"

	"gorm.io/gorm"
)

type KCRKScan_RFSS struct {
	SCNO     string    `gorm:"column:SCNO;primary_key;not null" json:"SCNO"`
	CLBH     string    `gorm:"column:CLBH;primary_key;not null" json:"CLBH"`
	Pack     int       `gorm:"column:Pack;primary_key;not null" json:"Pack"`
	Qty      float64   `gorm:"column:Qty" json:"Qty"`
	PrintS   int       `gorm:"column:PrintS" json:"PrintS"`
	Barcode  string    `gorm:"column:Barcode" json:"Barcode"`
	LotNO    string    `gorm:"column:LotNO" json:"LotNO"`
	Memo_RY  string    `gorm:"column:Memo_RY" json:"Memo_RY"`
	Box      string    `gorm:"column:Box" json:"Box"`
	USERID   string    `gorm:"column:USERID" json:"USERID"`
	USERDATE time.Time `gorm:"column:USERDATE" json:"USERDATE"`
	XXCC     string    `gorm:"column:XXCC" json:"XXCC"`
}

func (KCRKScan_RFSS) TableName() string {
	return "KCRKScan_RFSS"
}

func (m *KCRKScan_RFSS) BeforeCreate(*gorm.DB) error {
	m.USERDATE = time.Now()
	return nil
}

func (m *KCRKScan_RFSS) BeforeUpdate(*gorm.DB) error {
	m.USERDATE = time.Now()
	return nil
}
