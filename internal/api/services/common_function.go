package services

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"tyxuan-web-printlabel-api/internal/pkg/config"
	"tyxuan-web-printlabel-api/internal/pkg/database"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
)

type CommonFunction struct {
	*BaseService
}

var (
	iYear  string
	iMonth string
	SCNO1  string
	NDate  time.Time
)

var Commonf = &CommonFunction{}

func (s *CommonFunction) GetDateInfo() error {
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

	// Query to retrieve current year, month, and date
	err = db.Raw(`SELECT YEAR(getdate()) AS FY`).Scan(&iYear).Error
	if err != nil {
		return err
	}
	err = db.Raw(`SELECT RIGHT('0' + CAST(MONTH(GETDATE()) AS VARCHAR(2)), 2) AS FM`).Scan(&iMonth).Error
	if err != nil {
		return err
	}
	err = db.Raw(`SELECT getdate() AS NDate`).Scan(&NDate).Error
	if err != nil {
		return err
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) CreateRF(CGNO string, UserID string) (SCNOO string, err error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return "", err
	}
	dbInstance, _ := db.DB()

	// Query to find the latest SCNO
	var SCNO string
	err = db.Raw(`SELECT TOP 1 SCNO FROM KCRKScan_RF WHERE SCNO LIKE ? ORDER BY SCNO DESC `, iYear+iMonth+"%").Scan(&SCNO).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Calculate the new SCNO
	if SCNO != "" {
		SCNO = SCNO[6:10]
		SCNOInt, err := strconv.Atoi(SCNO)
		if err != nil {
			return "", err
		}
		SCNOInt++
		SCNO = fmt.Sprintf("%04d", SCNOInt)
	} else {
		SCNO = "0001"
	}
	SCNO = iYear + iMonth + SCNO

	// Query to retrieve GSBH
	var GSBH string
	err = db.Raw(`SELECT gsbh FROM CGZL WHERE CGNO = ?`, CGNO).Scan(&GSBH).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	query := `
		INSERT INTO KCRKScan_RF 
		(SCNO, GSBH, CKBH, USERDATE, USERID, MEMO, FIFO, YN, LB)
		VALUES (?, ?, ?, ?, ?, ?, ?, '1', '02')
	`

	// Insert a new record into KCRKScan_RF
	err = db.Exec(query, SCNO, GSBH, GSBH, NDate.Format("2006/01/02"), UserID, NDate.Format("2006-01-02"), strings.Split(NDate.Format("2006-01-02"), "-")[1]).Error
	if err != nil {
		return "", err
	}

	dbInstance.Close()
	if err != nil {
		return "", err
	}
	return SCNO, nil
}

func (s *CommonFunction) CreateRFS(SCNO string, CGNO string, CLBH string, Qty string, PackQty string, UserID string) error {
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
		INSERT INTO KCRKScan_RFS 
		(SCNO, CLBH, Qty, PackQty, USERDATE, USERID, YN, CGNO, SFL)
		VALUES (?, ?, ?, ?, ?, ?, '1', ?, 'TM')
	`

	err = db.Exec(query, SCNO, CLBH, Qty, PackQty, NDate.Format("2006/01/02"), UserID, CGNO).Error
	if err != nil {
		return err
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) CreateRFSS(SCNO string, CLBH string, Qty string, PackQty string, UserID string) error {
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

	var PackCount int
	RY := ""
	Article := ""
	var RemainQty, RFS_PackQty float64

	tmpQry := `
		SELECT * FROM KCRKScan_RFSS WHERE SCNO = ? AND CLBH = ?
	`

	rows, err := db.Raw(tmpQry, SCNO, CLBH).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return nil // Record already exists, do nothing
	}

	if qty, err := strconv.ParseFloat(Qty, 64); err == nil {
		if packQty, err := strconv.ParseFloat(PackQty, 64); err == nil {
			if qty > 0 && packQty > 0 {
				PackCount = 1
				RemainQty = qty
				RFS_PackQty = packQty

				for RemainQty > 0 {
					barcode := SCNO + CLBH + strconv.Itoa(PackCount)
					query := `
						INSERT INTO KCRKScan_RFSS
						(SCNO, CLBH, Pack, Qty, USERDATE, USERID, YN, Memo_RY, Memo_Article, barcode)
						VALUES (?, ?, ?, ?, ?, ?, '1', ?, ?, ?)
					`
					err := db.Exec(query, SCNO, CLBH, PackCount, math.Min(RFS_PackQty, RemainQty), NDate.Format("2006/01/02"), UserID, RY, Article, barcode).Error
					if err != nil {
						return err
					}
					RemainQty -= RFS_PackQty
					PackCount++
				}
			}
		}
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) CreateRFSSS(SCNO string, CLBH string, CGNO string, UserID string) error {
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

	var Pack int
	var CGNO1, ZLBH, XXCC string
	var TotalQty, tempQty, Qty float64

	rows, err := db.Raw("SELECT Pack,Qty FROM KCRKScan_RFSS WHERE SCNO = ? AND CLBH = ?", SCNO, CLBH).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&Pack, &TotalQty)
		if err != nil {
			return err
		}

		qry_zlbh, err := db.Raw(`
            SELECT CGZLSS.CGNO, CGZLSS.ZLBH, CGZLSS.XXCC, Round(CGZLSS.Qty-IsNull(CGRK.okQty,0.0),2) AS Qty
            FROM CGZLSS
            LEFT JOIN (
                SELECT CGNO, CLBH, ZLBH, XXCC, sum(Qty) AS okQty
                FROM KCRKScan_RFSSS
                LEFT JOIN KCRKScan_RF ON KCRKScan_RF.SCNO=KCRKScan_RFSSS.SCNO
                WHERE KCRKScan_RFSSS.CGNO= ? AND KCRKScan_RFSSS.CLBH= ? AND KCRKScan_RF.LB='02'
                GROUP BY CGNO, CLBH, ZLBH, XXCC
            ) CGRK ON CGRK.CGNO=CGZLSS.CGNO AND CGRK.ZLBH=CGZLSS.ZLBH AND CGRK.CLBH=CGZLSS.CLBH AND CGRK.XXCC=CGZLSS.XXCC
            WHERE CGZLSS.CGNO= ? AND CGZLSS.CLBH= ? AND Round(CGZLSS.Qty-IsNull(CGRK.okQty,0.0),2) > 0
        `, CGNO, CLBH, CGNO, CLBH).Rows()
		if err != nil {
			return err
		}
		defer qry_zlbh.Close()
		DetailQty := 0.0
		for qry_zlbh.Next() {
			err := qry_zlbh.Scan(&CGNO1, &ZLBH, &XXCC, &Qty)
			if err != nil {
				return err
			}
			DetailQty += Qty

			if TotalQty >= DetailQty {
				err := db.Exec(`
                    INSERT INTO KCRKScan_RFSSS (SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, USERDATE, USERID, YN)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '1')
                `, SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, NDate.Format("2006/01/02"), UserID).Error
				if err != nil {
					return err
				}
			} else {
				tempQty = DetailQty - TotalQty
				if tempQty < Qty {
					err := db.Exec(`
                        INSERT INTO KCRKScan_RFSSS (SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, USERDATE, USERID, YN)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '1')
					`, SCNO, CLBH, Pack, Qty-tempQty, ZLBH, XXCC, CGNO, NDate.Format("2006/01/02"), UserID).Error
					if err != nil {
						return err
					}
				}
				break
			}
		}
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) UpdateRFSSMemoRYMemoArticle(SCNO string, CLBH string) error {
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

	var Pack int
	rows, err := db.Raw("SELECT Pack FROM KCRKScan_RFSS WHERE SCNO = ? AND CLBH = ?", SCNO, CLBH).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&Pack)
		if err != nil {
			return err
		}

		err = db.Exec(`
		UPDATE KCRKScan_RFSS SET Memo_RY = KCRKScan_RFSSS.ZLBH, memo_Article = KCRKScan_RFSSS.SKU
		FROM (
			SELECT RFSSS.SCNO, RFSSS.CLBH, RFSSS.Pack,
				CAST(SUBSTRING((
						SELECT DISTINCT ',' + KCRKScan_RFSSS.ZLBH
						FROM KCRKScan_RFSSS
						WHERE KCRKScan_RFSSS.SCNO = RFSSS.SCNO AND KCRKScan_RFSSS.CLBH = RFSSS.CLBH AND KCRKScan_RFSSS.Pack = RFSSS.Pack
						FOR XML PATH('')
					), 2, 1000) AS VARCHAR(1000)) AS ZLBH,
				CAST(SUBSTRING((
						SELECT ',' + DDZL.Article + '/' + DDZL.BUYNO
						FROM KCRKScan_RFSSS
						LEFT JOIN DDZL ON DDZL.DDBH = KCRKScan_RFSSS.ZLBH
						WHERE KCRKScan_RFSSS.SCNO = RFSSS.SCNO AND KCRKScan_RFSSS.CLBH = RFSSS.CLBH AND KCRKScan_RFSSS.Pack = RFSSS.Pack
						GROUP BY DDZL.Article, DDZL.BUYNO
						FOR XML PATH('')
					), 2, 1000) AS VARCHAR(1000)) AS SKU
			FROM KCRKScan_RFSSS RFSSS
			WHERE SCNO = ? AND CLBH = ? AND Pack = ?
			GROUP BY RFSSS.SCNO, RFSSS.CLBH, RFSSS.Pack
		) KCRKScan_RFSSS
		WHERE KCRKScan_RFSS.SCNO = KCRKScan_RFSSS.SCNO AND KCRKScan_RFSS.CLBH = KCRKScan_RFSSS.CLBH AND KCRKScan_RFSS.Pack = KCRKScan_RFSSS.Pack
					`, SCNO, CLBH, Pack).Error
		if err != nil {
			return err
		}
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) UpdateRFSSMemoRYMemoArticleNoSize(SCNO string, CLBH string, CGNO string) error {
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
	err = db.Exec(`
		UPDATE KCRKScan_RFSS SET Memo_RY = Memo.ZLBH, memo_Article = Memo.SKU
		FROM (
			SELECT
				CAST(SUBSTRING((
						SELECT DISTINCT ',' + CGZLSS.ZLBH
						FROM CGZLSS
						WHERE CGZLSS.CGNO = CGZZ.CGNO AND CGZLSS.CLBH = CGZZ.CLBH 
						FOR XML PATH('')
					), 2, 1000) AS VARCHAR(1000)) AS ZLBH,
				CAST(SUBSTRING((
						SELECT ',' + DDZL.Article + '/' + DDZL.BUYNO
						FROM CGZLSS
						LEFT JOIN DDZL ON DDZL.DDBH = CGZLSS.ZLBH
						WHERE CGZLSS.CGNO = CGZZ.CGNO AND CGZLSS.CLBH = CGZZ.CLBH 
						GROUP BY DDZL.Article, DDZL.BUYNO
						FOR XML PATH('')
					), 2, 1000) AS VARCHAR(1000)) AS SKU
			FROM CGZLSS CGZZ
			where CGNO = ? and CLBH = ?
			) as Memo
		WHERE KCRKScan_RFSS.SCNO = ? AND KCRKScan_RFSS.CLBH = ?
	`, CGNO, CLBH, SCNO, CLBH).Error
	if err != nil {
		return err
	}
	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

// excel
func (s *CommonFunction) Import_Excel_Packqty(data []entities.Import_EX, UserID string) (mes string, e error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return "error connect", err
	}
	dbInstance, _ := db.DB()
	err = s.GetDateInfo()
	if err != nil {
		return "error date", err
	}
	var StringGCNOErorr string
	for i, r := range data {
		if strings.TrimSpace(r.CGNO) == "undefined" || strings.TrimSpace(r.CLBH) == "undefined" {
			continue
		}
		query := `
			SELECT ISNULL(MAX(KCRKScan_RFS.SCNO), '') AS SCNO,  ISNULL(CAST(MAX(Pack)as int) +1 ,'') as PACK,MEMO,DOCNO
			FROM KCRKScan_RFS
			LEFT JOIN KCRKScan_RF ON KCRKScan_RFS.SCNO = KCRKScan_RF.SCNO
			LEFT JOIN KCRKScan_RFSS ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
			WHERE ISNULL(KCRKScan_RFSS.CFMDel, '') = '' AND KCRKScan_RF.LB = '02'
			AND KCRKScan_RFS.CGNO = ? AND KCRKScan_RFS.CLBH = ?
			GROUP BY CGNO,MEMO,DOCNO
		`
		var SCNO string
		var MAXPACK string
		var Memo string
		var DOCNO string
		rows, err := db.Raw(query, r.CGNO, r.CLBH).Rows()
		if err != nil {
			return "error data", err
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&SCNO, &MAXPACK, &Memo, &DOCNO)
			if err != nil {
				return "error data", err
			}
		}

		var len int
		querycheck := `
			select COUNT(CGZLS.CGNO) as len from CGZLS left join CGZL on CGZL.CGNO=CGZLS.CGNO where CGZLS.CGNO = ? and CGZLS.CLBH = ? and CGZL.ZSBH = ?
		`
		err = db.Raw(querycheck, r.CGNO, r.CLBH, UserID).Scan(&len).Error
		if len < 1 {
			return "Đơn Hàng " + r.CGNO + " " + r.CLBH + " không tồn tại", err
		}
		if SCNO != "" {
			if r.MEMO != Memo {
				StringGCNOErorr += "STT" + strconv.Itoa(i+1) + " Đơn hàng " + r.CGNO + " - " + r.CLBH + " ngày " + Memo + " chưa Confirm !!!" + "\n"
				continue
			}
			if DOCNO != r.DOCNO {
				return "(" + strconv.Itoa(i+1) + ") Đơn hàng " + r.CGNO + " - " + r.CLBH + " ngày " + Memo + " DOCNO: " + DOCNO + " chưa Confirm !!!" + "\n", nil
			}
			//insert lại KCRKScan_RFSS
			barcode := SCNO + r.CLBH + MAXPACK
			insertQuery := `INSERT INTO KCRKScan_RFSS
							(SCNO, CLBH, Pack, Qty, USERDATE, USERID, YN, Memo_RY, Memo_Article,LotNO,barcode,Box)
							VALUES (?, ?, ?, ?, ?, ?, '1', ?, ?, ?,?,?)
							`
			err := db.Exec(insertQuery, SCNO, r.CLBH, MAXPACK, r.Qty, NDate.Format("2006/01/02"), UserID, "", "", r.LotNO, barcode, r.BoxNO).Error
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
			// Xóa dữ liệu từ K KCRKScan_RFSSS
			deleteQuery := `
				DELETE FROM KCRKScan_RFSSS WHERE SCNO = ? AND CLBH = ?
			`
			err = db.Exec(deleteQuery, SCNO, r.CLBH).Error
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
			//insert lại KCRKScan_RFSSS
			err = s.CreateRFSSS(SCNO, r.CLBH, r.CGNO, UserID)
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
			// Cập nhật RFSS_MemoRY và RFSS_MemoArticle
			err = s.UpdateRFSSMemoRYMemoArticleNoSize(SCNO, r.CLBH, r.CGNO)
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}

		} else {
			SCNO1, err = s.CreateRF(r.CGNO, UserID)
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}

			// Tạo mới RFS, RFSS, và RFSSS
			err = s.CreateRFS(SCNO1, r.CGNO, r.CLBH, r.Qty, r.Qty, UserID)
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
			//insert lại KCRKScan_RFSS
			barcode := SCNO1 + r.CLBH + "1"
			insertQuery := `INSERT INTO KCRKScan_RFSS
							(SCNO, CLBH, Pack, Qty, USERDATE, USERID, YN, Memo_RY, Memo_Article,LotNO,barcode,Box)
							VALUES (?, ?, ?, ?, ?, ?, '1', ?, ?, ?,?,?)
							`
			err := db.Exec(insertQuery, SCNO1, r.CLBH, 1, r.Qty, NDate.Format("2006/01/02"), UserID, "", "", r.LotNO, barcode, r.BoxNO).Error
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
			err = s.CreateRFSSS(SCNO1, r.CLBH, r.CGNO, UserID)
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
			// Cập nhật RFSS_MemoRY và RFSS_MemoArticle
			err = s.UpdateRFSSMemoRYMemoArticleNoSize(SCNO1, r.CLBH, r.CGNO)
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
		}
		// upload memo docon
		if SCNO != "" {
			queryupdate := `
				UPDATE KCRKScan_RFS
				SET Qty = KCRKScan_RFSS.Qty,  PackQty = KCRKScan_RFSS.PackQty
				FROM (		
					SELECT RFSS.SCNO, RFSS.CLBH, SUM(ISNULL(RFSS.Qty, 0)) AS Qty, max(Qty) as PackQty
					FROM KCRKScan_RFSS RFSS
					WHERE SCNO = ? AND CLBH = ?
					GROUP BY RFSS.SCNO, RFSS.CLBH
				) AS KCRKScan_RFSS
				WHERE KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
			`
			err = db.Exec(queryupdate, SCNO, r.CLBH).Error
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
			//
			updatequery := `UPDATE KCRKScan_RF
						SET DOCNO = ?, MEMO =?
						where  SCNO = ?`

			err = db.Exec(updatequery, r.DOCNO, r.MEMO, SCNO).Error
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
		} else {
			queryupdate := `
				UPDATE KCRKScan_RFS
				SET Qty = KCRKScan_RFSS.Qty,  PackQty = KCRKScan_RFSS.PackQty
				FROM (		
					SELECT RFSS.SCNO, RFSS.CLBH, SUM(ISNULL(RFSS.Qty, 0)) AS Qty, max(Qty) as PackQty
					FROM KCRKScan_RFSS RFSS
					WHERE SCNO = ? AND CLBH = ?
					GROUP BY RFSS.SCNO, RFSS.CLBH
				) AS KCRKScan_RFSS
				WHERE KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
			`
			err = db.Exec(queryupdate, SCNO1, r.CLBH).Error
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
			//
			updatequery := `UPDATE KCRKScan_RF
						SET DOCNO = ?, MEMO =?
						where  SCNO = ?`

			err = db.Exec(updatequery, r.DOCNO, r.MEMO, SCNO1).Error
			if err != nil {
				return "error data " + r.CGNO + " | " + r.CLBH + r.MEMO + strconv.Itoa(i), err
			}
		}
	}
	dbInstance.Close()
	if StringGCNOErorr == "" {
		StringGCNOErorr += "Nhập Excel thành công! Vui lòng kiểm tra lại!"
	}
	return StringGCNOErorr, nil
}

func (s *CommonFunction) SavePackqty(CGNO string, CLBH string, Qty string, PackQty string, UserID string) error {
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

	err = s.GetDateInfo()
	if err != nil {
		return err
	}

	query := `
        SELECT ISNULL(MAX(KCRKScan_RFS.SCNO), '') AS SCNO
        FROM KCRKScan_RFS
        LEFT JOIN KCRKScan_RF ON KCRKScan_RFS.SCNO = KCRKScan_RF.SCNO
        LEFT JOIN KCRKScan_RFSS ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
        WHERE ISNULL(KCRKScan_RFSS.CFMDel, '') = '' AND KCRKScan_RF.LB = '02'
        AND KCRKScan_RFS.CGNO = ? AND KCRKScan_RFS.CLBH = ?
    `

	var SCNO string
	err = db.Raw(query, CGNO, CLBH).Scan(&SCNO).Error
	if err != nil {
		return err
	}

	if SCNO != "" {
		// Cập nhật KCRKScan_RFS
		updateQuery := `
            UPDATE KCRKScan_RFS
            SET PackQty = ?, userdate = getdate()
            WHERE SCNO = ? AND CLBH = ?
        `
		err = db.Exec(updateQuery, PackQty, SCNO, CLBH).Error
		if err != nil {
			return err
		}

		// Xóa dữ liệu từ KCRKScan_RFSS và KCRKScan_RFSSS
		deleteQuery := `
            DELETE FROM KCRKScan_RFSS WHERE SCNO = ? AND CLBH = ?
            DELETE FROM KCRKScan_RFSSS WHERE SCNO = ? AND CLBH = ?
        `
		err = db.Exec(deleteQuery, SCNO, CLBH, SCNO, CLBH).Error
		if err != nil {
			return err
		}

		// Tạo mới RFSS và RFSSS
		err = s.CreateRFSS(SCNO, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}

		err = s.CreateRFSSS(SCNO, CLBH, CGNO, UserID)
		if err != nil {
			return err
		}

		// Cập nhật RFSS_MemoRY và RFSS_MemoArticle
		err = s.UpdateRFSSMemoRYMemoArticleNoSize(SCNO, CLBH, CGNO)
		if err != nil {
			return err
		}
	} else {
		// Tạo mới RF
		SCNO1, err = s.CreateRF(CGNO, UserID)
		if err != nil {
			return err
		}

		// Tạo mới RFS, RFSS, và RFSSS
		err = s.CreateRFS(SCNO1, CGNO, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}
		err = s.CreateRFSS(SCNO1, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}
		err = s.CreateRFSSS(SCNO1, CLBH, CGNO, UserID)
		if err != nil {
			return err
		}

		// Cập nhật RFSS_MemoRY và RFSS_MemoArticle
		err = s.UpdateRFSSMemoRYMemoArticleNoSize(SCNO1, CLBH, CGNO)
		if err != nil {
			return err
		}
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}
func (s *CommonFunction) Load(SCNO string, CGNO string, CLBH string, UserID string) error {
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

	err = s.GetDateInfo()
	if err != nil {
		return err
	}

	deleteQuery := `
        DELETE FROM KCRKScan_RFSSS WHERE SCNO = ? AND CLBH = ?
    `
	err = db.Exec(deleteQuery, SCNO, CLBH, SCNO, CLBH).Error
	if err != nil {
		return err
	}

	// Gọi hàm Create_RFSSS
	err = s.CreateRFSSS(SCNO, CLBH, CGNO, UserID)
	if err != nil {
		return err
	}

	// Gọi hàm Update_RFSS_MemoRY_MemoArticle
	err = s.UpdateRFSSMemoRYMemoArticleNoSize(SCNO, CLBH, CGNO)
	if err != nil {
		return err
	}

	// Update dữ liệu trong KCRKScan_RFS
	query := `
        UPDATE KCRKScan_RFS
        SET Qty = KCRKScan_RFSS.Qty
        FROM (
            SELECT RFSS.SCNO, RFSS.CLBH, SUM(ISNULL(RFSS.Qty, 0)) AS Qty
            FROM KCRKScan_RFSS RFSS
            WHERE SCNO = ? AND CLBH = ?
            GROUP BY RFSS.SCNO, RFSS.CLBH
        ) AS KCRKScan_RFSS
        WHERE KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
    `
	err = db.Exec(query, SCNO, CLBH).Error
	if err != nil {
		return err
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) CreateRFSSSZLBH(SCNO string, CLBH string, CGNO string, UserID string, ZLBH1 string) error {
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

	var Pack int
	var CGNO1, ZLBH, XXCC string
	var TotalQty, tempQty, Qty float64

	rows, err := db.Raw("SELECT Pack,Qty FROM KCRKScan_RFSS WHERE SCNO = ? AND CLBH = ?", SCNO, CLBH).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&Pack, &TotalQty)
		if err != nil {
			return err
		}

		qry_zlbh, err := db.Raw(`
            SELECT CGZLSS.CGNO, CGZLSS.ZLBH, CGZLSS.XXCC, Round(CGZLSS.Qty-IsNull(CGRK.okQty,0.0),2) AS Qty
            FROM CGZLSS
            LEFT JOIN (
                SELECT CGNO, CLBH, ZLBH, XXCC, sum(Qty) AS okQty
                FROM KCRKScan_RFSSS
                LEFT JOIN KCRKScan_RF ON KCRKScan_RF.SCNO=KCRKScan_RFSSS.SCNO
                WHERE KCRKScan_RFSSS.CGNO= ? AND KCRKScan_RFSSS.CLBH= ? and KCRKScan_RFSSS.ZLBH= ? AND KCRKScan_RF.LB='02'
                GROUP BY CGNO, CLBH, ZLBH, XXCC
            ) CGRK ON CGRK.CGNO=CGZLSS.CGNO AND CGRK.ZLBH=CGZLSS.ZLBH AND CGRK.CLBH=CGZLSS.CLBH AND CGRK.XXCC=CGZLSS.XXCC
            WHERE CGZLSS.CGNO= ? AND CGZLSS.CLBH= ? and CGZLSS.ZLBH= ? AND Round(CGZLSS.Qty-IsNull(CGRK.okQty,0.0),2) > 0
        `, CGNO, CLBH, ZLBH1, CGNO, CLBH, ZLBH1).Rows()
		if err != nil {
			return err
		}
		defer qry_zlbh.Close()
		DetailQty := 0.0
		for qry_zlbh.Next() {
			err := qry_zlbh.Scan(&CGNO1, &ZLBH, &XXCC, &Qty)
			if err != nil {
				return err
			}
			DetailQty += Qty

			if TotalQty >= DetailQty {
				err := db.Exec(`
                    INSERT INTO KCRKScan_RFSSS (SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, USERDATE, USERID, YN)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '1')
                `, SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, NDate.Format("2006/01/02"), UserID).Error
				if err != nil {
					return err
				}
			} else {
				tempQty = DetailQty - TotalQty
				if tempQty < Qty {
					err := db.Exec(`
                        INSERT INTO KCRKScan_RFSSS (SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, USERDATE, USERID, YN)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '1')
					`, SCNO, CLBH, Pack, Qty-tempQty, ZLBH, XXCC, CGNO, NDate.Format("2006/01/02"), UserID).Error
					if err != nil {
						return err
					}
				}
				break
			}
		}
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) CreateRFSSS_Size(SCNO string, CLBH string, CGNO string, UserID string, XXCC1 string) error {
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

	var Pack int
	var CGNO1, ZLBH, XXCC string
	var TotalQty, tempQty, Qty float64

	rows, err := db.Raw("SELECT Pack,Qty FROM KCRKScan_RFSS WHERE SCNO = ? AND CLBH = ?", SCNO, CLBH).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&Pack, &TotalQty)
		if err != nil {
			return err
		}

		qry_zlbh, err := db.Raw(`
            SELECT CGZLSS.CGNO, CGZLSS.ZLBH, CGZLSS.XXCC, Round(CGZLSS.Qty-IsNull(CGRK.okQty,0.0),2) AS Qty
            FROM CGZLSS
            LEFT JOIN (
                SELECT CGNO, CLBH, ZLBH, XXCC, sum(Qty) AS okQty
                FROM KCRKScan_RFSSS
                LEFT JOIN KCRKScan_RF ON KCRKScan_RF.SCNO=KCRKScan_RFSSS.SCNO
                WHERE KCRKScan_RFSSS.CGNO= ? AND KCRKScan_RFSSS.CLBH= ? and KCRKScan_RFSSS.XXCC= ? AND KCRKScan_RF.LB='02'
                GROUP BY CGNO, CLBH, ZLBH, XXCC
            ) CGRK ON CGRK.CGNO=CGZLSS.CGNO AND CGRK.ZLBH=CGZLSS.ZLBH AND CGRK.CLBH=CGZLSS.CLBH AND CGRK.XXCC=CGZLSS.XXCC
            WHERE CGZLSS.CGNO= ? AND CGZLSS.CLBH= ? and CGZLSS.XXCC= ? AND Round(CGZLSS.Qty-IsNull(CGRK.okQty,0.0),2) > 0
        `, CGNO, CLBH, XXCC1, CGNO, CLBH, XXCC1).Rows()
		if err != nil {
			return err
		}
		defer qry_zlbh.Close()
		DetailQty := 0.0
		for qry_zlbh.Next() {
			err := qry_zlbh.Scan(&CGNO1, &ZLBH, &XXCC, &Qty)
			if err != nil {
				return err
			}
			DetailQty += Qty

			if TotalQty >= DetailQty {
				err := db.Exec(`
                    INSERT INTO KCRKScan_RFSSS (SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, USERDATE, USERID, YN)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '1')
                `, SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, NDate.Format("2006/01/02"), UserID).Error
				if err != nil {
					return err
				}
			} else {
				tempQty = DetailQty - TotalQty
				if tempQty < Qty {
					err := db.Exec(`
                        INSERT INTO KCRKScan_RFSSS (SCNO, CLBH, Pack, Qty, ZLBH, XXCC, CGNO, USERDATE, USERID, YN)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '1')
					`, SCNO, CLBH, Pack, Qty-tempQty, ZLBH, XXCC, CGNO, NDate.Format("2006/01/02"), UserID).Error
					if err != nil {
						return err
					}
				}
				break
			}
		}
	}

	query := `
    UPDATE KCRKScan_RFSS
	SET XXCC = ? , userdate = getdate()
	WHERE SCNO = ? AND CLBH = ?
    `

	err = db.Raw(query, XXCC1, SCNO, CLBH).Scan(&SCNO).Error
	if err != nil {
		return err
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}
func (s *CommonFunction) SavePackqtyZLBH(CGNO string, CLBH string, Qty string, PackQty string, UserID string, ZLBH string) error {
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

	err = s.GetDateInfo()
	if err != nil {
		return err
	}

	query := `
        SELECT ISNULL(MAX(KCRKScan_RFS.SCNO), '') AS SCNO
        FROM KCRKScan_RFS
        LEFT JOIN KCRKScan_RF ON KCRKScan_RFS.SCNO = KCRKScan_RF.SCNO
        LEFT JOIN KCRKScan_RFSS ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
		left join KCRKScan_RFSSS on KCRKScan_RFSS.SCNO = KCRKScan_RFSSS.SCNO and KCRKScan_RFSS.CLBH = KCRKScan_RFSSS.CLBH and KCRKScan_RFSS.Pack = KCRKScan_RFSSS.Pack
        WHERE ISNULL(KCRKScan_RFSS.CFMDel, '') = '' AND KCRKScan_RF.LB = '02'
        AND KCRKScan_RFS.CGNO = ? AND KCRKScan_RFS.CLBH = ? and KCRKScan_RFSSS.ZLBH= ?
    `

	var SCNO string
	err = db.Raw(query, CGNO, CLBH, ZLBH).Scan(&SCNO).Error
	if err != nil {
		return err
	}

	if SCNO != "" {
		// Cập nhật KCRKScan_RFS
		updateQuery := `
            UPDATE KCRKScan_RFS
            SET PackQty = ?, userdate = getdate()
            WHERE SCNO = ? AND CLBH = ?
        `
		err = db.Exec(updateQuery, PackQty, SCNO, CLBH).Error
		if err != nil {
			return err
		}

		// Xóa dữ liệu từ KCRKScan_RFSS và KCRKScan_RFSSS
		deleteQuery := `
            DELETE FROM KCRKScan_RFSS WHERE SCNO = ? AND CLBH = ?
            DELETE FROM KCRKScan_RFSSS WHERE SCNO = ? AND CLBH = ?
        `
		err = db.Exec(deleteQuery, SCNO, CLBH, SCNO, CLBH).Error
		if err != nil {
			return err
		}

		// Tạo mới RFSS và RFSSS
		err = s.CreateRFSS(SCNO, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}

		err = s.CreateRFSSSZLBH(SCNO, CLBH, CGNO, UserID, ZLBH)
		if err != nil {
			return err
		}

		// Cập nhật RFSS_MemoRY và RFSS_MemoArticle
		err = s.UpdateRFSSMemoRYMemoArticle(SCNO, CLBH)
		if err != nil {
			return err
		}
	} else {
		// Tạo mới RF
		SCNO1, err = s.CreateRF(CGNO, UserID)
		if err != nil {
			return err
		}

		// Tạo mới RFS, RFSS, và RFSSS
		err = s.CreateRFS(SCNO1, CGNO, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}
		err = s.CreateRFSS(SCNO1, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}
		err = s.CreateRFSSSZLBH(SCNO1, CLBH, CGNO, UserID, ZLBH)
		if err != nil {
			return err
		}

		// Cập nhật RFSS_MemoRY và RFSS_MemoArticle
		err = s.UpdateRFSSMemoRYMemoArticle(SCNO1, CLBH)
		if err != nil {
			return err
		}
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) SavePackqtySize(CGNO string, CLBH string, Qty string, PackQty string, UserID string, XXCC string) error {
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

	err = s.GetDateInfo()
	if err != nil {
		return err
	}

	query := `
        SELECT ISNULL(MAX(KCRKScan_RFS.SCNO), '') AS SCNO
		FROM KCRKScan_RFS
		LEFT JOIN KCRKScan_RF ON KCRKScan_RFS.SCNO = KCRKScan_RF.SCNO
		LEFT JOIN KCRKScan_RFSS ON KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
		WHERE ISNULL(KCRKScan_RFSS.CFMDel, '') = '' AND KCRKScan_RF.LB = '02'
		AND KCRKScan_RFS.CGNO = ? AND KCRKScan_RFS.CLBH = ?  and XXCC = ?
    `

	var SCNO string
	err = db.Raw(query, CGNO, CLBH, XXCC).Scan(&SCNO).Error
	if err != nil {
		return err
	}

	if SCNO != "" {
		// Cập nhật KCRKScan_RFS
		updateQuery := `
            UPDATE KCRKScan_RFS
            SET PackQty = ?, userdate = getdate()
            WHERE SCNO = ? AND CLBH = ?
        `
		err = db.Exec(updateQuery, PackQty, SCNO, CLBH).Error
		if err != nil {
			return err
		}

		// Xóa dữ liệu từ KCRKScan_RFSS và KCRKScan_RFSSS
		deleteQuery := `
            DELETE FROM KCRKScan_RFSS WHERE SCNO = ? AND CLBH = ?
            DELETE FROM KCRKScan_RFSSS WHERE SCNO = ? AND CLBH = ?
        `
		err = db.Exec(deleteQuery, SCNO, CLBH, SCNO, CLBH).Error
		if err != nil {
			return err
		}

		// Tạo mới RFSS và RFSSS
		err = s.CreateRFSS(SCNO, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}

		err = s.CreateRFSSS_Size(SCNO, CLBH, CGNO, UserID, XXCC)
		if err != nil {
			return err
		}

		// Cập nhật RFSS_MemoRY và RFSS_MemoArticle
		err = s.UpdateRFSSMemoRYMemoArticle(SCNO, CLBH)
		if err != nil {
			return err
		}
	} else {
		// Tạo mới RF
		SCNO1, err = s.CreateRF(CGNO, UserID)
		if err != nil {
			return err
		}

		// Tạo mới RFS, RFSS, và RFSSS
		err = s.CreateRFS(SCNO1, CGNO, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}
		err = s.CreateRFSS(SCNO1, CLBH, Qty, PackQty, UserID)
		if err != nil {
			return err
		}
		err = s.CreateRFSSS_Size(SCNO1, CLBH, CGNO, UserID, XXCC)
		if err != nil {
			return err
		}

		// Cập nhật RFSS_MemoRY và RFSS_MemoArticle
		err = s.UpdateRFSSMemoRYMemoArticle(SCNO1, CLBH)
		if err != nil {
			return err
		}
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}
func (s *CommonFunction) LoadSize(SCNO string, CGNO string, CLBH string, UserID string, XXCC string) error {
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

	err = s.GetDateInfo()
	if err != nil {
		return err
	}

	deleteQuery := `
        DELETE FROM KCRKScan_RFSSS WHERE SCNO = ? AND CLBH = ?
    `
	err = db.Exec(deleteQuery, SCNO, CLBH, SCNO, CLBH).Error
	if err != nil {
		return err
	}

	// Gọi hàm Create_RFSSS
	err = s.CreateRFSSS_Size(SCNO, CLBH, CGNO, UserID, XXCC)
	if err != nil {
		return err
	}

	// Gọi hàm Update_RFSS_MemoRY_MemoArticle
	err = s.UpdateRFSSMemoRYMemoArticle(SCNO, CLBH)
	if err != nil {
		return err
	}

	// Update dữ liệu trong KCRKScan_RFS
	query := `
        UPDATE KCRKScan_RFS
        SET Qty = KCRKScan_RFSS.Qty
        FROM (
            SELECT RFSS.SCNO, RFSS.CLBH, SUM(ISNULL(RFSS.Qty, 0)) AS Qty
            FROM KCRKScan_RFSS RFSS
            WHERE SCNO = ? AND CLBH = ?
            GROUP BY RFSS.SCNO, RFSS.CLBH
        ) AS KCRKScan_RFSS
        WHERE KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
    `
	err = db.Exec(query, SCNO, CLBH).Error
	if err != nil {
		return err
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) LoadZLBH(SCNO string, CGNO string, CLBH string, UserID string, ZLBH string) error {
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

	err = s.GetDateInfo()
	if err != nil {
		return err
	}

	deleteQuery := `
        DELETE FROM KCRKScan_RFSSS WHERE SCNO = ? AND CLBH = ?
    `
	err = db.Exec(deleteQuery, SCNO, CLBH, SCNO, CLBH).Error
	if err != nil {
		return err
	}

	// Gọi hàm Create_RFSSS
	err = s.CreateRFSSSZLBH(SCNO, CLBH, CGNO, UserID, ZLBH)
	if err != nil {
		return err
	}

	// Gọi hàm Update_RFSS_MemoRY_MemoArticle
	err = s.UpdateRFSSMemoRYMemoArticle(SCNO, CLBH)
	if err != nil {
		return err
	}

	// Update dữ liệu trong KCRKScan_RFS
	query := `
        UPDATE KCRKScan_RFS
        SET Qty = KCRKScan_RFSS.Qty
        FROM (
            SELECT RFSS.SCNO, RFSS.CLBH, SUM(ISNULL(RFSS.Qty, 0)) AS Qty
            FROM KCRKScan_RFSS RFSS
            WHERE SCNO = ? AND CLBH = ?
            GROUP BY RFSS.SCNO, RFSS.CLBH
        ) AS KCRKScan_RFSS
        WHERE KCRKScan_RFS.SCNO = KCRKScan_RFSS.SCNO AND KCRKScan_RFS.CLBH = KCRKScan_RFSS.CLBH
    `
	err = db.Exec(query, SCNO, CLBH).Error
	if err != nil {
		return err
	}

	dbInstance.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CommonFunction) removeDiacritics(str string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: Nonspacing Marks
	}))
	str, _, _ = transform.String(t, str)
	return str
}
