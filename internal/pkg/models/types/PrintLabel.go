package types

type PurListPrint struct {
	SCNO    string
	CGNO    string
	CLBH    string
	YWPM    string
	ZLBH    string
	YQDate  string
	ZSYWJC  string
	PackQty float32
	Qty     float32
	DelQty  float32
	XXCC    string
	ARTICLE string
	KHPO    string
}

type LabelPrint struct {
	CLBH          string
	YWPM          string
	ZSYWJC        string
	CGNO          string
	XXCC          string
	Memo_RY       string
	Memo_Article  string
	Qty           string
	Pack          string
	LotNO         string
	Barcode       string
	FIFO          string
	Date_Received string
	Box           string
}

type DetailListPrint struct {
	SCNO    string
	CLBH    string
	Pack    int
	PrintS  int
	Qty     float32
	LotNO   string
	LotFile string
	Barcode string
}
