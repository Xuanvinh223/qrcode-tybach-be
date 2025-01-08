package request

import (
	"mime/multipart"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
)

type PurListRequest struct {
	ZSDH    string `form:"ZSDH" json:"ZSDH" binding:"required"`
	YQDate1 string `form:"YQDate1" json:"YQDate1" binding:"required"`
	YQDate2 string `form:"YQDate2" json:"YQDate2" binding:"required"`
	CGNO    string `form:"CGNO" json:"CGNO"`
	CLBH    string `form:"CLBH" json:"CLBH"`
	YWPM    string `form:"YWPM" json:"YWPM"`
}

type PurListRequestS struct {
	ZSDH    string `form:"ZSDH" json:"ZSDH" binding:"required"`
	YQDate1 string `form:"YQDate1" json:"YQDate1" binding:"required"`
	YQDate2 string `form:"YQDate2" json:"YQDate2" binding:"required"`
	CGNO    string `form:"CGNO" json:"CGNO"`
	CLBH    string `form:"CLBH" json:"CLBH"`
	YWPM    string `form:"YWPM" json:"YWPM"`
	ZLBH    string `form:"ZLBH" json:"ZLBH"`
}

type PurListRequestKHPO struct {
	ZSDH    string `form:"ZSDH" json:"ZSDH" binding:"required"`
	YQDate1 string `form:"YQDate1" json:"YQDate1" binding:"required"`
	YQDate2 string `form:"YQDate2" json:"YQDate2" binding:"required"`
	KHPO    string `form:"KHPO" json:"KHPO"`
	ARTICLE    string `form:"ARTICLE" json:"ARTICLE"`
	ZLBH    string `form:"ZLBH" json:"ZLBH"`
	YWPM    string `form:"YWPM" json:"YWPM"`
}

type PurListRequestSize struct {
	ZSDH    string `form:"ZSDH" json:"ZSDH" binding:"required"`
	YQDate1 string `form:"YQDate1" json:"YQDate1" binding:"required"`
	YQDate2 string `form:"YQDate2" json:"YQDate2" binding:"required"`
	CGNO    string `form:"CGNO" json:"CGNO"`
	CLBH    string `form:"CLBH" json:"CLBH"`
	YWPM    string `form:"YWPM" json:"YWPM"`
	XXCC    string `form:"XXCC" json:"XXCC"`
}

type DetailListRequest struct {
	CGNO string `form:"CGNO" json:"CGNO" binding:"required"`
	CLBH string `form:"CLBH" json:"CLBH" binding:"required"`
}

type DetailListRequestS struct {
	CGNO string `form:"CGNO" json:"CGNO" binding:"required"`
	CLBH string `form:"CLBH" json:"CLBH" binding:"required"`
	ZLBH string `form:"ZLBH" json:"ZLBH" binding:"required"`
}

type DetailListRequestSize struct {
	CGNO string `form:"CGNO" json:"CGNO" binding:"required"`
	CLBH string `form:"CLBH" json:"CLBH" binding:"required"`
	XXCC string `form:"XXCC" json:"XXCC" binding:"required"`
}

type PackQtyRequest struct {
	CGNO    string  `form:"CGNO" json:"CGNO" binding:"required"`
	CLBH    string  `form:"CLBH" json:"CLBH" binding:"required"`
	PackQty float64 `form:"PackQty" json:"PackQty" binding:"required"`
	USERID  string  `form:"USERID" json:"USERID" binding:"required"`
}

type Ip_ex_lPackQtyRequest struct {
	Data   []entities.Import_EX `form:"Data" json:"Data" binding:"required"`
	USERID string               `form:"USERID" json:"USERID" binding:"required"`
}

type PackQtyRequestS struct {
	CGNO    string  `form:"CGNO" json:"CGNO" binding:"required"`
	CLBH    string  `form:"CLBH" json:"CLBH" binding:"required"`
	ZLBH    string  `form:"ZLBH" json:"ZLBH" binding:"required"`
	PackQty float64 `form:"PackQty" json:"PackQty" binding:"required"`
	USERID  string  `form:"USERID" json:"USERID" binding:"required"`
}

type PackQtyRequestSize struct {
	CGNO    string  `form:"CGNO" json:"CGNO" binding:"required"`
	CLBH    string  `form:"CLBH" json:"CLBH" binding:"required"`
	XXCC    string  `form:"XXCC" json:"XXCC" binding:"required"`
	PackQty float64 `form:"PackQty" json:"PackQty" binding:"required"`
	USERID  string  `form:"USERID" json:"USERID" binding:"required"`
}

type PackRequest struct {
	SCNO   string `form:"SCNO" json:"SCNO" binding:"required"`
	CLBH   string `form:"CLBH" json:"CLBH" binding:"required"`
	USERID string `form:"USERID" json:"USERID" binding:"required"`
	QTYAD  int    `form:"QTYAD" json:"QTYAD"`
}

type PackRequestS struct {
	SCNO   string `form:"SCNO" json:"SCNO" binding:"required"`
	CLBH   string `form:"CLBH" json:"CLBH" binding:"required"`
	ZLBH   string `form:"ZLBH" json:"ZLBH" binding:"required"`
	USERID string `form:"USERID" json:"USERID" binding:"required"`
	QTYAD  int    `form:"QTYAD" json:"QTYAD"`
}

type PackRequestSize struct {
	SCNO   string `form:"SCNO" json:"SCNO" binding:"required"`
	CLBH   string `form:"CLBH" json:"CLBH" binding:"required"`
	USERID string `form:"USERID" json:"USERID" binding:"required"`
	QTYAD  int    `form:"QTYAD" json:"QTYAD"`
}

type PackRequest_SizeS struct {
	SCNO   string `form:"SCNO" json:"SCNO" binding:"required"`
	CLBH   string `form:"CLBH" json:"CLBH" binding:"required"`
	USERID string `form:"USERID" json:"USERID" binding:"required"`
	QTYAD  int    `form:"QTYAD" json:"QTYAD"`
	XXCC   string `form:"XXCC" json:"XXCC" binding:"required"`
}

type LoadRequest struct {
	CGNO   string                   `form:"CGNO" json:"CGNO" binding:"required"`
	DelQty float64                  `form:"DelQty" json:"DelQty" binding:"required"`
	DOCNO  string                   `form:"DOCNO" json:"DOCNO"`
	MEMO   string                   `form:"MEMO" json:"MEMO" binding:"required"`
	Data   []entities.KCRKScan_RFSS `form:"Data" json:"Data" binding:"required"`
}

type LoadRequestS struct {
	CGNO   string                   `form:"CGNO" json:"CGNO" binding:"required"`
	ZLBH   string                   `form:"ZLBH" json:"ZLBH" binding:"required"`
	DelQty float64                  `form:"DelQty" json:"DelQty" binding:"required"`
	DOCNO  string                   `form:"DOCNO" json:"DOCNO"`
	MEMO   string                   `form:"MEMO" json:"MEMO" binding:"required"`
	Data   []entities.KCRKScan_RFSS `form:"Data" json:"Data" binding:"required"`
}

type LoadRequestSize struct {
	CGNO   string                   `form:"CGNO" json:"CGNO" binding:"required"`
	XXCC   string                   `form:"XXCC" json:"XXCC" binding:"required"`
	DelQty float64                  `form:"DelQty" json:"DelQty" binding:"required"`
	DOCNO  string                   `form:"DOCNO" json:"DOCNO"`
	MEMO   string                   `form:"MEMO" json:"MEMO" binding:"required"`
	Data   []entities.KCRKScan_RFSS `form:"Data" json:"Data" binding:"required"`
}

type CFMRequest struct {
	SCNO   string `form:"SCNO" json:"SCNO" binding:"required"`
	CLBH   string `form:"CLBH" json:"CLBH" binding:"required"`
	USERID string `form:"USERID" json:"USERID" binding:"required"`
}

type UploadLotFileRequest struct {
	ZSDH    string               `form:"ZSDH" json:"ZSDH" binding:"required"`
	LotNO   string               `form:"LotNO" json:"LotNO" binding:"required"`
	LotFile multipart.FileHeader `form:"LotFile" json:"LotFile" binding:"required"`
}
