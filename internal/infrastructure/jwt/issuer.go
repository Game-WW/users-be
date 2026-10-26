package jwt

import (
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/Game-WW/users-be/internal/core/application/commands"
)

type Issuer struct {
	secret []byte
}

var _ commands.AccessTokenIssuer = (*Issuer)(nil)

func NewIssuer(secret string) *Issuer {
	return &Issuer{secret: []byte(secret)}
}

type claims struct {
	UserID string `json:"user_id"`
	jwtlib.RegisteredClaims
}

func (i *Issuer) Issue(userID string, expiresAt time.Time) (string, error) {
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims{
		UserID: userID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
		},
	})

	signed, err := token.SignedString(i.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}
