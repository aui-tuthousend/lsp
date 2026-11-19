package auth

import "context"

type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Name   string `json:"name"`
}

type claimsKey struct{}

func ContextWithClaims(ctx context.Context, claims UserClaims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

func ClaimsFromContext(ctx context.Context) (UserClaims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(UserClaims)
	return claims, ok
}
