package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"tyxuan-web-printlabel-api/internal/pkg/models/response"
	"tyxuan-web-printlabel-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func JWTAuth(ctx *gin.Context) {
	authorizationHeader := ctx.Request.Header.Get(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "authorization header is not provide")
		ctx.Abort()
		return
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, "invalid authorization header format")
		ctx.Abort()
		return
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		err := fmt.Sprintf("unsupported authorization type %s", authorizationType)
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, err)
		ctx.Abort()
		return
	}

	accessToken := fields[1]
	payload, err := jwt.ParseToken(accessToken)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, err.Error())
		ctx.Abort()
		return
	}

	ctx.Set(authorizationPayloadKey, payload)
	ctx.Next()
}
