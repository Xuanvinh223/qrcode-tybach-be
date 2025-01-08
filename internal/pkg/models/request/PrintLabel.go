package request

import (
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
)

type PurListPrintRequest struct {
	ZSDH    string `form:"ZSDH" json:"ZSDH" binding:"required"`
	YQDate1 string `form:"YQDate1" json:"YQDate1" binding:"required"`
	YQDate2 string `form:"YQDate2" json:"YQDate2" binding:"required"`
	CGNO    string `form:"CGNO" json:"CGNO"`
	CLBH    string `form:"CLBH" json:"CLBH"`
	YWPM    string `form:"YWPM" json:"YWPM"`
	ZLBH    string `form:"ZLBH" json:"ZLBH"`
	XXCC    string `form:"XXCC" json:"XXCC"`
}

type DetailListPrintRequest struct {
	SCNO string `form:"SCNO" json:"SCNO" binding:"required"`
	CLBH string `form:"CLBH" json:"CLBH" binding:"required"`
}

type DetailListPrintRequest1 struct {
	Data []DetailListPrintRequest `form:"Data" json:"Data" binding:"required"`
}

type LabelPrintRequest struct {
	Data []entities.KCRKScan_RFSS `form:"Data" json:"Data" binding:"required"`
}
