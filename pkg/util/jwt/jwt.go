package jwt

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt"
	"time"
)

var ErrGenerateTokenFailed = errors.New(500, "GENERATE_TOKEN_FAILED", "generate token failed")

type CustomClaimsConfiguration struct {
	UserID      uint64
	NickName    string
	AuthorityId uint64
	BelongTo    uint64
	SecretKey   string
	Timeout     int64
}
type CustomClaims struct {
	ID          uint64
	NickName    string
	AuthorityId uint64
	BelongTo    uint64
	jwt.StandardClaims
}

func GenerateToken(config *CustomClaimsConfiguration) (string, error) {
	t := time.Now().Unix()
	claims := CustomClaims{
		ID:          config.UserID,
		NickName:    config.NickName,
		BelongTo:    config.BelongTo,
		AuthorityId: config.AuthorityId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: t + config.Timeout,
			Issuer:    config.NickName,
			NotBefore: t,
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := tkn.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", ErrGenerateTokenFailed
	}
	return signedString, nil
}
