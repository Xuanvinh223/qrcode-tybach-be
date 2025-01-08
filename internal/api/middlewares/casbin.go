package middlewares

import (
	_ "embed"
	"net/http"
	"strings"

	"tyxuan-web-printlabel-api/internal/pkg/database"
	"tyxuan-web-printlabel-api/internal/pkg/models/response"
	"tyxuan-web-printlabel-api/pkg/jwt"
	util2 "tyxuan-web-printlabel-api/pkg/util"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
)

//go:embed rbac_model.conf
var rbacModelConf string

func CasbinAuth(ctx *gin.Context) {
	db := database.GetDB()

	adapter, _ := gormadapter.NewAdapterByDB(db)
	rc, err := model.NewModelFromString(rbacModelConf)
	if err != nil {
		ctx.Abort()
		return
	}
	e, _ := casbin.NewEnforcer(rc, adapter)
	e.AddFunction("ParamsMatch", ParamsMatchFunc)
	e.AddFunction("ParamsActMatch", ParamsActMatchFunc)
	_ = e.LoadPolicy()

	// current request
	obj := ctx.Request.URL.RequestURI()
	act := ctx.Request.Method
	user := ctx.MustGet("authorization_payload").(*jwt.UserClaims)
	if len(user.RoleID) == 0 {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "Authorization Exception")
		ctx.Abort()
		return
	}

	var flag = false
	for _, sub := range user.RoleID {
		// 判斷策略中是否存在
		subStr := gconv.String(sub)
		if ok, _ := e.Enforce(subStr, obj, act); ok {
			flag = true
			break
		}
	}
	if !flag {
		response.FailWithDetailed(ctx, http.StatusForbidden, nil, "User does not have permission to perform this action")
		ctx.Abort()
		return
	}
	ctx.Next()
}

// ParamsMatchFunc 自定義規則函數
func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)
	key1 := strings.Split(name1, "?")[0]
	// 剝離路徑後再使用 casbin 的 keyMatch2
	return util.KeyMatch2(key1, name2), nil
}

// ParamsActMatchFunc 自定義規則函數
func ParamsActMatchFunc(args ...interface{}) (interface{}, error) {
	rAct := args[0].(string)
	pAct := args[1].(string)
	pActArr := strings.Split(pAct, ",")
	return util2.InAnySlice[string](pActArr, rAct), nil
}
