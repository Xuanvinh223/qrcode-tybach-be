package types

type PurCFM struct {
	SCNO    string
	CGNO    string
	CLBH    string
	YWPM    string
	YQDate  string
	ZSYWJC  string
	ZLBH    string
	DOCNO   string
	MEMO    string
	Status  string
	PackQty float32
	Qty     float32
	DelQty  float32
	KHPO    string
	ARTICLE string
}

type DetailCFM struct {
	SCNO    string
	CLBH    string
	Pack    int
	PrintS  int
	Qty     float32
	LotNO   string
	LotFile string
	Barcode string
}
