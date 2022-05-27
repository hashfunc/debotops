package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Kind string

type Claims struct {
	jwt.RegisteredClaims

	Kind Kind `json:"kind"`
}

const (
	KindAccessToken  Kind = "ACCESS"
	KindRefreshToken Kind = "REFRESH"
)
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = time.Hour
)

func NewToken(key string, kind Kind) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expires(kind),
		},
		Kind: kind,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func expires(kind Kind) *jwt.NumericDate {
	switch kind {
	case KindAccessToken:
		return jwt.NewNumericDate(time.Now().Add(AccessTokenDuration))
	case KindRefreshToken:
		return jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration))
	default:
		return jwt.NewNumericDate(time.Now())
	}
}
