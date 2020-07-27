package util

import (
	"time"
	"tree-hole/config"

	"github.com/dgrijalva/jwt-go"
)

type jwtClaims struct {
	jwt.StandardClaims
	ID string `json:"_id"`
}

func NewJWTToken(idHex string) (string, time.Time, error) {
	expireTime := time.Now().Add(time.Duration(config.Config.JWT.TokenExpire) * time.Minute)
	claims := jwtClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   "tree-hole",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime.Unix(),
		},
		ID: idHex,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(config.Config.JWT.Secret))
	if err != nil {
		return "", time.Now(), err
	}
	return tokenSigned, expireTime, nil
}
