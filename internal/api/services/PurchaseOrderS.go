package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"tyxuan-web-printlabel-api/internal/pkg/config"
	"tyxuan-web-printlabel-api/internal/pkg/database"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"
)

type PursService struct {
	*BaseService
	*CommonFunction
}

var Purs = &PursService{}

func (s *PursService) GetPurListS(requestParams request.PurListRequestS) ([]types.PurListS, error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return []types.PurListS{}, err
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`
	WITH PackData AS (
		SELECT 
			RFS.CLBH, 
			RFS.CGNO, 
			RFSS.Memo_RY AS ZLBH, 
			MAX(ISNULL(RFS.PackQty, 0)) AS PackQty
		FROM KCRKScan_RFSS RFSS
		LEFT JOIN KCRKScan_RFS RFS ON RFSS.SCNO = RFS.SCNO AND RFS.CLBH = RFSS.CLBH
		LEFT JOIN KCRKScan_RF RF ON RFSS.SCNO = RF.SCNO
		WHERE RF.LB = '02' AND ISNULL(RFSS.CFMDel, '') = ''
		GROUP BY RFS.CLBH, RFS.CGNO, RFSS.Memo_RY
	),
	KCRKScanSS AS (
		SELECT 
			RFS.CLBH, 
			RFS.CGNO, 
			RFSS.Memo_RY AS ZLBH, 
			SUM(RFSS.qty) AS TotalQty, 
			MAX(RFSS.CFMDel) AS CFMDel, 
			MAX(ISNULL(RFS.PackQty, 0)) AS PackQty
		FROM KCRKScan_RFSS RFSS
		LEFT JOIN KCRKScan_RFS RFS ON RFSS.SCNO = RFS.SCNO AND RFS.CLBH = RFSS.CLBH
		LEFT JOIN KCRKScan_RF RF ON RFSS.SCNO = RF.SCNO
		WHERE RF.LB = '02'
		GROUP BY RFS.CLBH, RFS.CGNO, RFSS.Memo_RY
	),
	MainQuery AS (
		SELECT 
			CGZLSS.CGNO, 
			CGZLSS.CLBH, 
			CLZL.YWPM, 
			CONVERT(varchar, CGZLSS.YQDate, 111) AS YQDate, 
			ZSZL.ZSDH, 
			ZSZL.zsywjc, 
			ISNULL(PackData.PackQty, 0) AS PackQty, 
			SUM(CGZLSS.qty) AS Qty, 
			ISNULL(KCRKScanSS.TotalQty, 0) AS DelQty, 
			CGZLSS.ZLBH
		FROM CGZLSS
		LEFT JOIN CGZL ON CGZL.CGNO = CGZLSS.CGNO
		LEFT JOIN ZSZL ON CGZL.ZSBH = ZSZL.ZSDH
		LEFT JOIN CLZL ON CGZLSS.CLBH = CLZL.CLDH
		LEFT JOIN PackData ON 
			PackData.CGNO = CGZLSS.CGNO AND 
			PackData.CLBH = CGZLSS.CLBH AND 
			PackData.ZLBH = CGZLSS.ZLBH
		LEFT JOIN KCRKScanSS ON 
			KCRKScanSS.CGNO = CGZLSS.CGNO AND 
			KCRKScanSS.CLBH = CGZLSS.CLBH AND 
			KCRKScanSS.ZLBH = CGZLSS.ZLBH
		WHERE 
			CGZL.ZSBH = '%s' 
			AND CGZLSS.XXCC <> 'ZZZZZZ'
			AND CGZLSS.YQDate BETWEEN '%s' AND '%s'
			AND CGZLSS.CGNO LIKE '%s%%'
			AND CGZLSS.CLBH LIKE '%s%%'
			AND CLZL.YWPM LIKE '%%%s%%'
			AND CGZLSS.ZLBH LIKE '%%%s%%'
		GROUP BY 
			CGZLSS.CGNO, CGZLSS.CLBH, CLZL.YWPM, CGZLSS.ZLBH, 
			CONVERT(varchar, CGZLSS.YQDate, 111), ZSZL.ZSDH, ZSZL.ZSYWJC, 
			PackData.PackQty, KCRKScanSS.TotalQty
	)
	SELECT *
	FROM MainQuery
	WHERE 
		DelQty < Qty OR 
		(DelQty >= Qty AND ISNULL(PackQty, '') <> '')
	ORDER BY YQDate DESC, CGNO, CLBH, ZLBH; `,
		requestParams.ZSDH, requestParams.YQDate1, requestParams.YQDate2, requestParams.CGNO, requestParams.CLBH, requestParams.YWPM, requestParams.ZLBH)

	var result []types.PurListS
	err = db.Raw(query).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PursService) SetPackQtyS(requestParams request.PackQtyRequestS) error {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return err
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`
	select sum(cgzlss.Qty) Qty,max(isnull(KCRKScan_RFSS.qty,0)) as DelQty 
	from cgzlss left join cgzl on cgzl.cgno=cgzlss.cgno 
	left join ZSZL on CGZL.ZSBH = ZSZL.ZSDH left join CLZL on cgzlss.CLBH = CLZL.CLDH 
	left join (
	Select KCRKScan_RFSS.CLBH,KCRKScan_RFS.CGNO, KCRKScan_RFSS.Memo_RY ZLBH, SUM(KCRKScan_RFSS.qty) as qty, max(KCRKScan_RFSS.CFMDel) as CFMDel
    from KCRKScan_RFSS left join KCRKScan_RF on KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO 
    left join KCRKScan_RFS on KCRKScan_RFS.SCNO=KCRKScan_RFSS.SCNO 
    and KCRKScan_RFS.CLBH=KCRKScan_RFSS.CLBH where KCRKScan_RF.LB='02' and isnull(KCRKScan_RFSS.CFMDel,'')<>''
    group by KCRKScan_RFSS.CLBH,KCRKScan_RFS.CGNO, KCRKScan_RFSS.Memo_RY) KCRKScan_RFSS 
	on KCRKScan_RFSS.CGNO=cgzlss.cgno and KCRKScan_RFSS.CLBH=cgzlss.CLBH and KCRKScan_RFSS.ZLBH=cgzlss.ZLBH
	where CGZLSS.CGNO='%s' and CGZLSS.CLBH='%s' and CGZLSS.ZLBH='%s' `, requestParams.CGNO, requestParams.CLBH, requestParams.ZLBH)

	var result types.PackQty
	err = db.Raw(query).Scan(&result).Error
	var ReQty float64
	ReQty = result.Qty - result.DelQty

	err = s.SavePackqtyZLBH(requestParams.CGNO, requestParams.CLBH, strconv.FormatFloat(ReQty, 'f', -1, 64), strconv.FormatFloat(requestParams.PackQty, 'f', -1, 64),
		requestParams.USERID, requestParams.ZLBH)
	if err != nil {
		return err
	}
	// requestURL := fmt.Sprintf("http://192.168.23.58:8001/?CFM=0&CGNO=%s&CLBH=%s&Qty=%.2f&PackQty=%.2f&UserID=%s",
	// 	requestParams.CGNO, requestParams.CLBH, ReQty, requestParams.PackQty, requestParams.USERID)
	// _, err = http.Get(requestURL)

	dbInstance.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *PursService) GetDetailListS(requestParams request.DetailListRequestS) ([]types.DetailListS, error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return []types.DetailListS{}, err
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`
	select KCRKScan_RFSS.SCNO, KCRKScan_RFSS.CLBH, KCRKScan_RFSS.Memo_RY ZLBH, KCRKScan_RFSS.Pack, KCRKScan_RFSS.Qty, KCRKScan_RFSS.LotNO,
	 	   KCRKScan_RFS.CGNO, zszl_lot.LotFile, zszl_lot.ZSDH, DOCNO, MEMO
	from KCRKScan_RFSS 
	left join KCRKScan_RF on KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO 
	left join KCRKScan_RFS on KCRKScan_RFSS.SCNO = KCRKScan_RFS.SCNO and KCRKScan_RFSS.CLBH = KCRKScan_RFS.CLBH 
	left join
	(SELECT zszl_lot.zsdh,zszl_lot.LotNO,CAST(substring (( select case when isnull(lot.LotFile,'')<>'' then ', ' + isnull(lot.LotFile,'') else '' end
	FROM zszl_lot lot  WHERE lot.zsdh=zszl_lot.zsdh and lot.LotNO=zszl_lot.LotNO
	 for XML Path ('')),2,1000) as varchar(1000)) as LotFile
	FROM zszl_lot
	group by zszl_lot.zsdh,zszl_lot.LotNO) AS zszl_lot on zszl_lot.zsdh=KCRKScan_RFSS.USERID and zszl_lot.LotNO=KCRKScan_RFSS.LotNO 
	where KCRKScan_RF.LB='02' and isnull(KCRKScan_RFSS.CFMDel,'')='' 
	and KCRKScan_RFS.CGNO='%s' and KCRKScan_RFS.CLBH='%s' and KCRKScan_RFSS.Memo_RY='%s'
	order by KCRKScan_RFSS.Pack asc `, requestParams.CGNO, requestParams.CLBH, requestParams.ZLBH)

	var result []types.DetailListS
	err = db.Raw(query).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PursService) AddPackS(SCNO string, CLBH string, ZLBH string, USERID string, QTYAD int) error {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"

	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return err
	}

	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	for i := 0; i < QTYAD; i++ {
		query := fmt.Sprintf(`
		select Max(KCRKScan_RFSS.Pack) as Pack
		from KCRKScan_RFSS    
		where KCRKScan_RFSS.SCNO='%s' and KCRKScan_RFSS.CLBH='%s' `, SCNO, CLBH)

		var Pack int
		err = db.Raw(query).Scan(&Pack).Error
		Pack = Pack + 1

		object := entities.KCRKScan_RFSS{
			SCNO:    SCNO,
			CLBH:    CLBH,
			Pack:    Pack,
			Qty:     0,
			PrintS:  0,
			Barcode: SCNO + CLBH + strconv.Itoa(Pack),
			LotNO:   "",
			Memo_RY: ZLBH,
			USERID:  USERID,
		}
		_ = db.Create(&object)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PursService) DelPackS(SCNO string, CLBH string, USERID string, QTYAD int) error {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"

	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return err
	}

	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	for i := 0; i < QTYAD; i++ {
		query := fmt.Sprintf(`
		select Max(KCRKScan_RFSS.Pack) as Pack
		from KCRKScan_RFSS    
		where KCRKScan_RFSS.SCNO='%s' and KCRKScan_RFSS.CLBH='%s' `, SCNO, CLBH)

		var Pack int
		err = db.Raw(query).Scan(&Pack).Error

		deleteQuery := `
        DELETE FROM KCRKScan_RFSSS WHERE SCNO = ? AND CLBH = ? AND Pack = ?
    	`
		err = db.Exec(deleteQuery, SCNO, CLBH, Pack).Error
		if err != nil {
			return err
		}

		where := entities.KCRKScan_RFSS{
			SCNO: SCNO,
			CLBH: CLBH,
			Pack: Pack,
		}

		_ = db.Where(&where).Delete(&entities.KCRKScan_RFSS{})

		if Pack == 1 {
			deleteQuery = `
			DELETE FROM KCRKScan_RFS WHERE SCNO = ? AND CLBH = ? 
			`
			err = db.Exec(deleteQuery, SCNO, CLBH).Error
			if err != nil {
				return err
			}
			deleteQuery = `
			DELETE FROM KCRKScan_RF WHERE SCNO = ?
			`
			err = db.Exec(deleteQuery, SCNO).Error
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PursService) LoadingQtyS(requestParams request.LoadRequestS) error {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return err
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`
	select max(CGZLSS.Qty)-max(isnull(KCRKScan_RFSSS.Qty,0)) as AllQty
	from KCRKScan_RFSS    
	left join KCRKScan_RF on KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO    
	left join KCRKScan_RFS on KCRKScan_RFSS.SCNO = KCRKScan_RFS.SCNO and KCRKScan_RFSS.CLBH = KCRKScan_RFS.CLBH
	left join (select KCRKScan_RFSS.SCNO, CLBH, sum(isnull(qty,0)) qty from KCRKScan_RFSS
	left join KCRKScan_RF on KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO
	where KCRKScan_RF.LB='02' and isnull(KCRKScan_RFSS.CFMDel,'')<>''
	group by KCRKScan_RFSS.SCNO, CLBH) as KCRKScan_RFSSS
	on KCRKScan_RFSS.SCNO = KCRKScan_RFSSS.SCNO and KCRKScan_RFSS.CLBH = KCRKScan_RFSSS.CLBH
	left join (select CGNO, CLBH, ZLBH, sum (CGZLSS.Qty) Qty from CGZLSS
	group by CGNO, CLBH, ZLBH) as CGZLSS 
	on CGZLSS.CGNO = KCRKScan_RFS.CGNO and CGZLSS.CLBH = KCRKScan_RFSS.CLBH and CGZLSS.ZLBH=KCRKScan_RFSS.Memo_RY
	where KCRKScan_RF.LB='02' and isnull(KCRKScan_RFSS.CFMDel,'')=''
	and KCRKScan_RFSS.SCNO='%s' and KCRKScan_RFSS.CLBH='%s' `, requestParams.Data[0].SCNO, requestParams.Data[0].CLBH)

	var AllQty float64
	err = db.Raw(query).Scan(&AllQty).Error

	if requestParams.DelQty <= AllQty {
		updateQuery := `
            UPDATE KCRKScan_RF
            SET DOCNO = ?, MEMO = ?, FIFO = ?
            WHERE SCNO = ?
        `
		err = db.Exec(updateQuery, requestParams.DOCNO, requestParams.MEMO, strings.Split(requestParams.MEMO, "-")[1], requestParams.Data[0].SCNO).Error
		if err != nil {
			return err
		}

		for _, object := range requestParams.Data {
			float, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", object.Qty), 64)
			value := entities.KCRKScan_RFSS{
				Qty:   float,
				LotNO: object.LotNO,
			}

			where := entities.KCRKScan_RFSS{
				SCNO: object.SCNO,
				CLBH: object.CLBH,
				Pack: object.Pack,
			}

			if err := db.Model(&where).Updates(&value); err != nil {
				continue
			}
		}
	}
	err = s.LoadZLBH(requestParams.Data[0].SCNO, requestParams.CGNO, requestParams.Data[0].CLBH, requestParams.Data[0].USERID, requestParams.ZLBH)
	if err != nil {
		return err
	}
	// requestURL := fmt.Sprintf("http://192.168.23.58:8001/?CFM=1&SCNO=%s&CLBH=%s&CGNO=%s&CFMDel=0&UserID=%s",
	// 	requestParams.Data[0].SCNO, requestParams.Data[0].CLBH, requestParams.CGNO, requestParams.Data[0].USERID)
	// _, err = http.Get(requestURL)

	dbInstance.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *PursService) UploadLotFileS(requestParams request.UploadLotFileRequest, LotFile *multipart.FileHeader) error {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return err
	}
	dbInstance, _ := db.DB()

	// fileReader, err := LotFile.Open()
	// if err != nil {
	// 	return err
	// }

	// defer fileReader.Close()

	// if err != nil {
	// 	return err
	// }

	// // Create a buffer to store the request body
	// var buf bytes.Buffer

	// // Create a new multipart writer with the buffer
	// w := multipart.NewWriter(&buf)

	// // Add a file to the request
	// //file, err := os.Open("./uploadfile/" + LotFile.Filename)
	// file, err := LotFile.Open()
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// // Create a new form field
	// fw, err := w.CreateFormFile("LotFile", LotFile.Filename)
	// if err != nil {
	// 	return err
	// }

	// // Copy the contents of the file to the form field
	// if _, err := io.Copy(fw, file); err != nil {
	// 	return err
	// }

	// // Close the multipart writer to finalize the request
	// w.Close()

	// // Send the request
	// req, err := http.NewRequest("POST", "http://192.168.23.58/QR/LotFile/uploadfile.php", &buf)
	// if err != nil {
	// 	return err
	// }
	// req.Header.Set("Content-Type", w.FormDataContentType())

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// Mở file nguồn
	file, err := LotFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Tạo tên tệp mới mà không có dấu
	newFileName := s.removeDiacritics(LotFile.Filename)

	// Đảm bảo rằng tên tệp không chứa khoảng trắng hoặc ký tự đặc biệt
	newFileName = strings.Map(func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		return r
	}, newFileName)

	// Tạo file đích
	dst, err := os.Create("./uploadfile/" + newFileName)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy nội dung từ file nguồn vào file đích
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	object := entities.Zszl_lot{
		ZSDH:    requestParams.ZSDH,
		LotNO:   requestParams.LotNO,
		LotFile: newFileName,
	}
	_ = db.Create(&object)
	if err := db.Save(&object); err != nil {
		return err.Error
	}

	dbInstance.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *PursService) CFMS(requestParams request.CFMRequest) error {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return err
	}
	dbInstance, _ := db.DB()

	query := `
		update KCRKScan_RFSS set CFMDel = ?, CFMDelDate=getdate()
        WHERE SCNO = ? AND CLBH = ?
    `
	err = db.Exec(query, requestParams.USERID, requestParams.SCNO, requestParams.CLBH).Error
	if err != nil {
		return err
	}

	dbInstance.Close()
	if err != nil {
		return err
	}

	return nil
}
