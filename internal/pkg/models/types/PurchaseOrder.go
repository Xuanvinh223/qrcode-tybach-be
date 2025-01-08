package types

type PurList struct {
	CGNO    string
	CLBH    string
	YWPM    string
	YQDate  string
	ZSDH    string
	ZSYWJC  string
	PackQty float32
	Qty     float32
	DelQty  float32
}
type PurListMes struct {
	Mes string
}

type PurListS struct {
	CGNO    string
	CLBH    string
	YWPM    string
	YQDate  string
	ZSDH    string
	ZSYWJC  string
	ZLBH    string
	PackQty float32
	Qty     float32
	DelQty  float32
}

type PurList_Nosize_KHPO struct {
	CGNO    string
	KHPO    string
	ARTICLE string
	CLBH    string
	YWPM    string
	YQDate  string
	ZSDH    string
	ZSYWJC  string
	PackQty float32
	Qty     float32
	DelQty  float32
	ZLBH    string
}

type PurListSize struct {
	CGNO    string
	CLBH    string
	YWPM    string
	YQDate  string
	ZSDH    string
	ZSYWJC  string
	XXCC    string
	PackQty float32
	Qty     float32
	DelQty  float32
}

type PackQty struct {
	Qty    float64
	DelQty float64
}

type DetailList struct {
	SCNO    string
	CLBH    string
	Pack    int
	Qty     float32
	CGNO    string
	LotNO   string
	LotFile string
	ZSDH    string
	DOCNO   string
	MEMO    string
	Box     string
}

type DetailListS struct {
	SCNO    string
	CLBH    string
	Pack    int
	Qty     float32
	CGNO    string
	ZLBH    string
	LotNO   string
	LotFile string
	ZSDH    string
	DOCNO   string
	MEMO    string
}

type DetailListSize struct {
	SCNO    string
	CLBH    string
	Pack    int
	Qty     float32
	XXCC    string
	CGNO    string
	LotNO   string
	LotFile string
	ZSDH    string
	DOCNO   string
	MEMO    string
}
