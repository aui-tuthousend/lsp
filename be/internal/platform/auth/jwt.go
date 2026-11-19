package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret []byte
	expiry time.Duration
}

type jwtClaims struct {
	UserClaims
	jwt.RegisteredClaims
}

func NewJWTProvider(secret string, expiry time.Duration) *JWTProvider {
	return &JWTProvider{secret: []byte(secret), expiry: expiry}
}

func (p *JWTProvider) GenerateToken(claims UserClaims) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		UserClaims:       claims,
		RegisteredClaims: jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(p.expiry)), Subject: claims.UserID},
	})
	return token.SignedString(p.secret)
}

func (p *JWTProvider) VerifyToken(tokenString string) (UserClaims, error) {
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return p.secret, nil
	})
	if err != nil || !token.Valid {
		return UserClaims{}, errors.New("invalid token")
	}
	return claims.UserClaims, nil
}
