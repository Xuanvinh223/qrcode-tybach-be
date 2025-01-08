package entities

type Import_EX struct {
	CGNO  string `gorm:"column:CGNO;primary_key;not null" json:"CGNO"`
	CLBH  string `gorm:"column:CLBH;primary_key;not null" json:"CLBH"`
	BoxNO string `gorm:"column:BoxNO;primary_key;not null" json:"BoxNO"`
	LotNO string `gorm:"column:LotNO" json:"LotNO"`
	Qty   string `gorm:"column:Qty" json:"Qty"`
	DOCNO string `gorm:"column:DOCNO" json:"DOCNO"`
	MEMO  string `gorm:"column:MEMO" json:"MEMO"`
}
