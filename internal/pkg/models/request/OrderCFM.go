package request

import (
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
)

type CFMAllRequest struct {
	Data []entities.KCRKScan_RFSS `form:"Data" json:"Data" binding:"required"`
}

type PurCFMRequest struct {
	ZSDH    string `form:"ZSDH" json:"ZSDH" binding:"required"`
	YQDate1 string `form:"YQDate1" json:"YQDate1" binding:"required"`
	YQDate2 string `form:"YQDate2" json:"YQDate2" binding:"required"`
	CGNO    string `form:"CGNO" json:"CGNO"`
	CLBH    string `form:"CLBH" json:"CLBH"`
	YWPM    string `form:"YWPM" json:"YWPM"`
	ZLBH    string `form:"ZLBH" json:"ZLBH"`
	CFM     int    `form:"CFM" json:"CFM" binding:"required"`
}

type DetailCFMRequest struct {
	SCNO string `form:"SCNO" json:"SCNO" binding:"required"`
	CLBH string `form:"CLBH" json:"CLBH" binding:"required"`
}
