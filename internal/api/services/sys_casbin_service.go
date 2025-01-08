package services

import (
	"strings"

	"tyxuan-web-printlabel-api/internal/pkg/database"
	"tyxuan-web-printlabel-api/internal/pkg/models/entities"
	"tyxuan-web-printlabel-api/internal/pkg/models/request"
	"tyxuan-web-printlabel-api/internal/pkg/models/types"
	"tyxuan-web-printlabel-api/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
)

type SysCasbinService struct {
	*BaseService
}

var SysCasbin = &SysCasbinService{}

func (s *SysCasbinService) AddRoutes(requestParams request.RouteRequest, routerInstance *gin.Engine) ([]entities.SysCasbin, error) {
	var (
		casbinRule     = make([]entities.SysCasbin, 0)
		routeList      = make([]types.RouteInfo, 0)
		newPathList    = make([]string, 0)
		casbinRuleList = make([]entities.SysCasbin, 0)
		db             = database.GetDB()
	)

	for _, route := range routerInstance.Routes() {
		var _routeInfo = types.RouteInfo{}
		if err := copier.Copy(&_routeInfo, route); err != nil {
			continue
		}
		routerArr := strings.Split(route.Path, "/")[1:]
		if len(routerArr) < 2 {
			_routeInfo.NewPath = "/" + strings.Join(routerArr, "/")
		} else {
			routerArr = strings.Split(route.Path, "/")[1:len(routerArr)]
			_routeInfo.NewPath = "/" + strings.Join(routerArr, "/") + "/*"
		}
		if strings.Contains(_routeInfo.NewPath, requestParams.RType) {
			routeList = append(routeList, _routeInfo)
			newPathList = append(newPathList, _routeInfo.NewPath)
		}
	}
	newPathList = lo.Uniq[string](newPathList)
	_routeList := make([]types.RouteInfo, len(newPathList))

	for i := 0; i < len(newPathList); i++ {
		_routeList[i].NewPath = newPathList[i]
		for j := 0; j < len(routeList); j++ {
			if routeList[j].NewPath == newPathList[i] {
				_routeList[i].MethodList = append(_routeList[i].MethodList, routeList[j].Method)
			}
		}
	}
	for i := 0; i < len(_routeList); i++ {
		_routeList[i].MethodList = lo.Uniq[string](_routeList[i].MethodList)
		casbinRule = append(casbinRule, entities.SysCasbin{
			Ptype: "p",
			V0:    gconv.String(requestParams.RoleId),
			V1:    _routeList[i].NewPath,
			V2:    strings.Join(_routeList[i].MethodList, ","),
		})
	}

	// 存入到 casbin_rule 中
	for _, rule := range casbinRule {
		var casbinRuleInfo entities.SysCasbin
		if err := db.Where("v0 = ? and v1 = ?", requestParams.RoleId, rule.V1).First(&casbinRuleInfo).Error; err != nil {
			casbinRuleList = append(casbinRuleList, rule)
			continue
		}
		// 存在
		flag := true
		ruleV2 := strings.Split(rule.V2, ",")
		casbinRuleInfoV2 := strings.Split(casbinRuleInfo.V2, ",")
		for _, s2 := range ruleV2 {
			if !util.InAnySlice(casbinRuleInfoV2, s2) {
				// 不存在
				flag = false
				break
			}
		}
		if !flag {
			// 刪除紀錄
			if err := db.Where("v0 = ? and v1 = ?", requestParams.RoleId, rule.V1).Delete(&casbinRuleInfo).Error; err != nil {
				continue
			}
			casbinRuleList = append(casbinRuleList, rule)
		}
	}
	if len(casbinRuleList) > 0 {
		if err := db.Model(entities.SysCasbin{}).Create(casbinRuleList).Error; err != nil {
			return nil, err
		}
	}
	return casbinRuleList, nil
}
