package jwt

import (
	"time"
)

var (
	TokenExpireDuration = time.Hour * 2
	MySecret            = []byte("123456")
)

// func GenToken(uid string) (string, error) {
// 	claims := MyClaims{
// 		Uid: uid,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 30 * time.Hour)), // 过期时间
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),                          // 签发时间
// 			NotBefore: jwt.NewNumericDate(time.Now()),                          // 生效时间
// 		},
// 	}
// 	// 使用HS256签名算法
// 	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	s, err := t.SignedString(MySecret)

// 	return s, err
// }

// func ParseJwt(tokenstring string) (*MyClaims, error) {
// 	t, err := jwt.ParseWithClaims(tokenstring, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(MySecret), nil
// 	})

// 	if claims, ok := t.Claims.(*MyClaims); ok && t.Valid {
// 		return claims, nil
// 	} else {
// 		return nil, err
// 	}
// }
