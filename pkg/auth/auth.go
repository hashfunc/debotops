package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Kind string

const (
	KindAccess  Kind = "ACCESS"
	KindRefresh Kind = "REFRESH"
)

type Claims struct {
	jwt.RegisteredClaims

	Kind Kind `json:"kind"`
}

type Auth struct {
	Kind    Kind
	Expires time.Time
	Token   string
}

func (token *Auth) String() string {
	return token.Token
}

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = time.Hour
)

func NewAuth(key string, kind Kind) (*Auth, error) {
	expires := newExpires(kind)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
		Kind: kind,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}

	auth := &Auth{
		Kind:    kind,
		Expires: expires,
		Token:   signed,
	}

	return auth, nil
}

func newExpires(kind Kind) time.Time {
	now := time.Now()
	switch kind {
	case KindAccess:
		return now.Add(AccessTokenDuration)
	case KindRefresh:
		return now.Add(RefreshTokenDuration)
	default:
		return now
	}
}
