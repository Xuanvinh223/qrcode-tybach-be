package services

import (
	"fmt"

	//"strings"

	"tyxuan-web-printlabel-api/internal/pkg/config"
	"tyxuan-web-printlabel-api/internal/pkg/database"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"
)

type PurPrintService struct {
	*BaseService
}

var PurP = &PurPrintService{}

func (s *PurPrintService) GetPurListPrint(requestParams request.PurListPrintRequest) ([]types.PurListPrint, error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return []types.PurListPrint{}, err
	}
	dbInstance, _ := db.DB()
	var where string
	if requestParams.ZSDH == "VA05" {
		where = "and DDZL.KHPO like '%%" + requestParams.CGNO + "%%'"
	} else {
		where = "and KCRKScan_RFS.CGNO like '%%" + requestParams.CGNO + "%%'"
	}

	query := fmt.Sprintf(`	
		select KCRKScan_RFS.SCNO, KCRKScan_RFS.CGNO, DDZL.KHPO, DDZL.ARTICLE, KCRKScan_RFS.CLBH, CLZL.YWPM, KCRKScan_RFS.Qty, KCRKScan_RFS.PackQty, ZSZL.ZSYWJC, CONVERT(varchar, CGZLS.YQDate, 111) as YQDate,
			isnull(KCRKScan_RFSS.qty,0) as DelQty, KCRKScan.ZLBH,CASE  WHEN XXCC = 'ZZZZZZ' THEN 'No Size' ELSE XXCC END AS XXCC
		from KCRKScan_RFS   
		left join KCRKScan_RF on KCRKScan_RFS.SCNO = KCRKScan_RF.SCNO    
		left join ( Select SCNO, CLBH, MAX(Memo_RY) ZLBH, MAX(ISNULL(XXCC,'')) XXCC
							from KCRKScan_RFSS
							group by SCNO, CLBH
					)KCRKScan on KCRKScan.SCNO=KCRKScan_RFS.SCNO and KCRKScan.CLBH=KCRKScan_RFS.CLBH
		left join( Select 	KCRKScan_RFSS.SCNO, KCRKScan_RFSS.CLBH, SUM(isnull(KCRKScan_RFSS.qty,0)) as qty, max(CFMDel) as CFMDel
							from KCRKScan_RFSS 
							left join KCRKScan_RF on KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO    
							where KCRKScan_RF.LB='02' AND CFMDel is not null 
							group by KCRKScan_RFSS.SCNO, CLBH) KCRKScan_RFSS on KCRKScan_RFS.SCNO=KCRKScan_RFSS.SCNO and KCRKScan_RFS.CLBH=KCRKScan_RFSS.CLBH 
		left join CGZL on KCRKScan_RFS.CGNO = CGZL.CGNO    
		left join CGZLS on CGZLS.CGNO = CGZL.CGNO and KCRKScan_RFS.CLBH=CGZLS.CLBH   
		left join ZSZL on CGZL.ZSBH = ZSZL.ZSDH    
		left join CLZL on KCRKScan_RFS.CLBH = CLZL.CLDH 
		left join DDZL on KCRKScan.ZLBH = DDZL.DDBH 
		WHERE KCRKScan_RF.LB='02' and ZSZL.zsdh = '%s' and convert(smalldatetime,convert(varchar,CGZLS.YQDate,111)) between '%s' and '%s'
		 and KCRKScan_RFS.CLBH like '%%%s%%' and CLZL.YWPM like '%%%s%%' and KCRKScan.ZLBH like'%%%s%%' and XXCC like '%%%s%%' %s
		order by KCRKScan_RFS.CGNO asc  
	`, requestParams.ZSDH, requestParams.YQDate1, requestParams.YQDate2, requestParams.CLBH, requestParams.YWPM, requestParams.ZLBH, requestParams.XXCC, where)

	var result []types.PurListPrint
	err = db.Raw(query).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return []types.PurListPrint{}, err
	}

	return result, nil
}

func (s *PurPrintService) GetDetailListPrint(requestParams request.DetailListPrintRequest1) ([]types.DetailListPrint, error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return []types.DetailListPrint{}, err
	}
	dbInstance, _ := db.DB()

	var where string
	//SCNOList := strings.Split(SCNO, "|")
	//CLBHList := strings.Split(CLBH, "|")

	where = fmt.Sprintf(` (KCRKScan_RFSS.SCNO='%s' and KCRKScan_RFSS.CLBH='%s') `, requestParams.Data[0].SCNO, requestParams.Data[0].CLBH)

	for i := 1; i < len(requestParams.Data); i++ {
		where = fmt.Sprintf(`%s or (KCRKScan_RFSS.SCNO='%s' and KCRKScan_RFSS.CLBH='%s') `, where, requestParams.Data[i].SCNO, requestParams.Data[i].CLBH)
	}

	query := fmt.Sprintf(`
	select SCNO, CLBH, Pack, PrintS, Qty, KCRKScan_RFSS.LotNO, zszl_lot.LotFile, Barcode
	from KCRKScan_RFSS 
	left join
	(SELECT zszl_lot.zsdh,zszl_lot.LotNO,CAST(substring (( select case when isnull(lot.LotFile,'')<>'' then ', ' + isnull(lot.LotFile,'') else '' end
	FROM zszl_lot lot  WHERE lot.zsdh=zszl_lot.zsdh and lot.LotNO=zszl_lot.LotNO
	 for XML Path ('')),2,1000) as varchar(1000)) as LotFile
	FROM zszl_lot
	group by zszl_lot.zsdh,zszl_lot.LotNO) AS zszl_lot on zszl_lot.zsdh=KCRKScan_RFSS.USERID and zszl_lot.LotNO=KCRKScan_RFSS.LotNO
	where %s
	order by Pack asc `, where)

	var result []types.DetailListPrint
	err = db.Raw(query).Scan(&result).Error
	dbInstance.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PurPrintService) GetLabelPrint(requestParams request.LabelPrintRequest) ([]types.LabelPrint, error) {
	configuration := config.GetConfig()
	configuration.Database.Driver = "sqlserver"
	configuration.Database.Host = "192.168.40.9"
	configuration.Database.Username = "tyxuan"
	configuration.Database.Password = "jack"
	configuration.Database.Dbname = "TB_ERP"
	configuration.Database.Port = "1433"
	db, err := database.CreateDatabaseConnection(configuration)
	if err != nil {
		return []types.LabelPrint{}, err
	}
	dbInstance, _ := db.DB()
	var result []types.LabelPrint
	for _, object := range requestParams.Data {

		value := entities.KCRKScan_RFSS{
			PrintS: object.PrintS + 1,
		}

		where := entities.KCRKScan_RFSS{
			SCNO: object.SCNO,
			CLBH: object.CLBH,
			Pack: object.Pack,
		}

		var label types.LabelPrint

		query := fmt.Sprintf(`
			select KCRKScan_RFSS.CLBH, YWPM, ZSYWJC, KCRKScan_RFS.CGNO, isnull(KCRKScan_RFSS.Memo_RY,'') as Memo_RY, isnull(KCRKScan_RFSS.Memo_Article,'') as Memo_Article, 
			CONVERT(varchar,KCRKScan_RFSS.Qty)+'  '+isnull(CLZL.DWBH,'') as Qty, CONVERT(varchar, KCRKScan_RFSS.Pack,101)+' of ' +CONVERT(varchar,TotalPack,101) as Pack, 
			LotNO, Barcode, isnull(KCRKScan_RF.FIFO,CONVERT(varchar,month(getdate()))) as fifo, isnull(KCRKScan_RF.Memo,'')+' | '+isnull(KCRKScan_RF.DOCNO,'') Date_Received, 
			KCRKScan_RFSS.Box,isnull( CASE  WHEN KCRKScan_RFSS.XXCC = 'ZZZZZZ' THEN '' ELSE 'Size: '+KCRKScan_RFSS.XXCC END,'') as XXCC
			from KCRKScan_RFSS 
			left join KCRKScan_RF on KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO 
			left join KCRKScan_RFS on KCRKScan_RFSS.SCNO = KCRKScan_RFS.SCNO and KCRKScan_RFSS.CLBH = KCRKScan_RFS.CLBH 
			left join (	select SCNO,CLBH,count (Pack) as TotalPack   
						from KCRKScan_RFSS   
						group by SCNO,CLBH) KCRKScan  on KCRKScan.SCNO = KCRKScan_RFS.SCNO and KCRKScan.CLBH = KCRKScan_RFS.CLBH
			left join CGZL on KCRKScan_RFS.CGNO = CGZL.CGNO 
			left join ZSZL on CGZL.ZSBH = ZSZL.ZSDH 
			left join CLZL on KCRKScan_RFSS.CLBH = CLZL.CLDH 
			where KCRKScan_RFSS.SCNO ='%s' and KCRKScan_RFSS.CLBH ='%s' and KCRKScan_RFSS.Pack ='%d'
			order by Pack asc`, object.SCNO, object.CLBH, object.Pack)

		err = db.Raw(query).Scan(&label).Error
		if err != nil {
			continue
		}
		result = append(result, label)

		if err := db.Model(&where).Updates(&value); err != nil {
			continue
		}
	}

	dbInstance.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
