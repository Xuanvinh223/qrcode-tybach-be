package services

import (
	"fmt"

	"tyxuan-web-printlabel-api/internal/pkg/config"
	"tyxuan-web-printlabel-api/internal/pkg/database"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"
)

type PurCFMService struct {
	*BaseService
	*CommonFunction
}

var PurCFM = &PurCFMService{}

func (s *PurCFMService) GetPurCFM(requestParams request.PurCFMRequest) ([]types.PurCFM, error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyx"
	configuration.Database.Password = "tyx"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return []types.PurCFM{}, err
	}
	dbInstance, _ := db.DB()
	var sql = ""
	if requestParams.CFM == 1 {
		sql = "and KCRKScan_RFSS.CFMDel is null"
	} else {
		sql = "and KCRKScan_RFSS.CFMDel is not null"
	}
	var where string
	if requestParams.ZSDH == "VA05" {
		where = "and DDZL.KHPO like '%%" + requestParams.CGNO + "%%'"
	} else {
		where = "and KCRKScan_RFS.CGNO like '%%" + requestParams.CGNO + "%%'"
	}

	query := fmt.Sprintf(`
	select 	KCRKScan_RFS.SCNO, DDZL.KHPO,DDZL.ARTICLE, KCRKScan_RFS.CGNO, KCRKScan_RFS.CLBH, CLZL.YWPM, KCRKScan_RFS.Qty, KCRKScan_RFS.PackQty, ZSZL.ZSYWJC, CONVERT(varchar, CGZLS.YQDate, 111) as YQDate,
			isnull(KCRKScan_RFSS.qty,0) as DelQty, KCRKScan_RF.DOCNO, KCRKScan_RF.MEMO, case when KCRKScan_RFSS.CFMDel is null then 'Not Cofirm' else 'Confirm' end as Status,
			CAST(SUBSTRING((
						SELECT DISTINCT ',' + RFSSS.ZLBH
						FROM KCRKScan_RFSSS RFSSS
						WHERE KCRKScan_RFSS.SCNO = RFSSS.SCNO AND KCRKScan_RFSS.CLBH = RFSSS.CLBH
						FOR XML PATH('')
					), 2, 1000) AS VARCHAR(1000)) AS ZLBH
    from KCRKScan_RFS   
	left join (SELECT SCNO,CLBH,Memo_RY FROM KCRKScan_RFSS GROUP BY SCNO,CLBH,Memo_RY) RFSS on RFSS.SCNO=KCRKScan_RFS.SCNO AND KCRKScan_RFS.CLBH=RFSS.CLBH
	left join DDZL on DDZL.DDBH = RFSS.Memo_RY
    left join KCRKScan_RF on KCRKScan_RFS.SCNO = KCRKScan_RF.SCNO    
    left join( Select 	KCRKScan_RFSS.SCNO, KCRKScan_RFSS.CLBH, SUM(isnull(KCRKScan_RFSS.qty,0)) as qty, max(CFMDel) as CFMDel
                        from KCRKScan_RFSS 
                        left join KCRKScan_RF on KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO    
                        where KCRKScan_RF.LB='02'
                        group by KCRKScan_RFSS.SCNO, CLBH) KCRKScan_RFSS on KCRKScan_RFS.SCNO=KCRKScan_RFSS.SCNO and KCRKScan_RFS.CLBH=KCRKScan_RFSS.CLBH 
    left join CGZL on KCRKScan_RFS.CGNO = CGZL.CGNO    
    left join CGZLS on CGZLS.CGNO = CGZL.CGNO and KCRKScan_RFS.CLBH=CGZLS.CLBH   
    left join ZSZL on CGZL.ZSBH = ZSZL.ZSDH    
    left join CLZL on KCRKScan_RFS.CLBH = CLZL.CLDH 
	WHERE KCRKScan_RF.LB='02' and ZSZL.zsdh = '%s' and convert(smalldatetime,convert(varchar,CGZLS.YQDate,111)) between '%s' and '%s'
	%s and KCRKScan_RFS.CLBH like '%%%s%%' and CLZL.YWPM like '%%%s%%' and isnull(CAST(SUBSTRING((
		SELECT DISTINCT ',' + RFSSS.ZLBH
		FROM KCRKScan_RFSSS RFSSS
		WHERE KCRKScan_RFSS.SCNO = RFSSS.SCNO AND KCRKScan_RFSS.CLBH = RFSSS.CLBH
		FOR XML PATH('')
	), 2, 1000) AS VARCHAR(1000)),'') like '%%%s%%' and KCRKScan_RFSS.SCNO is not null %s 
	order by KCRKScan_RFS.CGNO asc 
`,
		requestParams.ZSDH, requestParams.YQDate1, requestParams.YQDate2, where, requestParams.CLBH, requestParams.YWPM, requestParams.ZLBH, sql)

	var result []types.PurCFM
	err = db.Raw(query).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return []types.PurCFM{}, err
	}

	return result, nil
}

func (s *PurCFMService) GetDetailCFM(requestParams request.DetailCFMRequest) ([]types.DetailCFM, error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyx"
	configuration.Database.Password = "tyx"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return []types.DetailCFM{}, err
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`
	select SCNO, CLBH, Pack, PrintS, Qty, KCRKScan_RFSS.LotNO, zszl_lot.LotFile, Barcode
	from KCRKScan_RFSS 
	left join
	(SELECT zszl_lot.zsdh,zszl_lot.LotNO,CAST(substring (( select case when isnull(lot.LotFile,'')<>'' then ', ' + isnull(lot.LotFile,'') else '' end
	FROM zszl_lot lot  WHERE lot.zsdh=zszl_lot.zsdh and lot.LotNO=zszl_lot.LotNO
	 for XML Path ('')),2,1000) as varchar(1000)) as LotFile
	FROM zszl_lot
	group by zszl_lot.zsdh,zszl_lot.LotNO) AS zszl_lot on zszl_lot.zsdh=KCRKScan_RFSS.USERID and zszl_lot.LotNO=KCRKScan_RFSS.LotNO
	where KCRKScan_RFSS.SCNO='%s' and KCRKScan_RFSS.CLBH='%s'
	order by Pack asc `, requestParams.SCNO, requestParams.CLBH)

	var result []types.DetailCFM
	err = db.Raw(query).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PurCFMService) CFMALL(requestParams request.CFMAllRequest) error {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyx"
	configuration.Database.Password = "tyx"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return err
	}
	dbInstance, _ := db.DB()

	for _, object := range requestParams.Data {
		query := `
		update KCRKScan_RFSS set CFMDel = ?, CFMDelDate=getdate()
        WHERE SCNO = ? AND CLBH = ?
    	`
		if err = db.Exec(query, object.USERID, object.SCNO, object.CLBH).Error; err != nil {
			continue
		}
	}

	if err != nil {
		return err
	}

	dbInstance.Close()
	if err != nil {
		return err
	}

	return nil
}
