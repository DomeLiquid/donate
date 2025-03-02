package jwt

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// type Jwt interface {
// 	GenToken(uid string) (string, error)
// 	ParseJwt(tokenstring string) (*MyClaims, error)
// }

// var _ Jwt = (*jwtImpl)(nil)

type JwtConfig struct {
	SecretKey          string `mapstructure:"secret_key"`
	TokenExpireSeconds int64  `mapstructure:"token_expire_seconds"`
}

var (
	once sync.Once
	_jwt = &jwtImpl{}
)

func Init(config *JwtConfig) {
	once.Do(func() {
		_jwt = &jwtImpl{config: config}
	})
}

type MyClaims struct {
	Uid                  string `json:"uid"`
	jwt.RegisteredClaims        // v5版本新加的方法
}

type jwtImpl struct {
	config *JwtConfig
}

func GenToken(uid string) (string, error) {
	claims := MyClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(_jwt.config.TokenExpireSeconds * int64(time.Second)))), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                         // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                                         // 生效时间
		},
	}
	// 使用HS256签名算法
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(_jwt.config.SecretKey))

	return s, err
}

func ParseJwt(tokenstring string) (*MyClaims, error) {
	t, err := jwt.ParseWithClaims(tokenstring, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(_jwt.config.SecretKey), nil
	})

	if claims, ok := t.Claims.(*MyClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
