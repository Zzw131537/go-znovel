package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("Zhouzw")

type Claims struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Authority int    `json:"authority"`
	jwt.StandardClaims
}

// 签发
func GenerateToken(id uint, userName string, authority int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		ID:        id,
		Username:  userName,
		Authority: authority,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Zhouzw-novel",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err

}

// 验证用户toekn
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
