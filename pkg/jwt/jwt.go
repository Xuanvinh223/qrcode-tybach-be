package jwt

import (
	"time"

	"tyxuan-web-printlabel-api/internal/pkg/config"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID uint64   `json:"userId"`
	RoleID []uint64 `json:"roleId"`
}

func GenerateToken(claims *UserClaims) string {
	secret := []byte(config.GetConfig().Server.Secret)
	effectTime := 2 * time.Hour

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(effectTime))

	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		panic(err)
	}
	return sign
}

func ParseToken(tokenString string) (*UserClaims, error) {
	secret := []byte(config.GetConfig().Server.Secret)

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("unable to parse claims")
	}
	return claims, nil
}
