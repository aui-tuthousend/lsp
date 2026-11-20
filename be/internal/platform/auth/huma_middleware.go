package auth

import (
	"encoding/json"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

func NewJWTAuthMiddleware(provider *JWTProvider) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if isPublic(ctx.Operation()) {
			next(ctx)
			return
		}
		token := bearerToken(ctx)
		if token == "" {
			writeUnauthorized(ctx, "missing bearer token")
			return
		}
		claims, err := provider.VerifyToken(token)
		if err != nil {
			writeUnauthorized(ctx, "invalid bearer token")
			return
		}
		next(huma.WithContext(ctx, ContextWithClaims(ctx.Context(), claims)))
	}
}

// bearerToken mengambil token dari header Authorization, atau dari query
// `?token=` sebagai fallback untuk koneksi streaming (SSE) yang tidak bisa
// mengirim header custom dari EventSource browser.
func bearerToken(ctx huma.Context) string {
	if header := ctx.Header("Authorization"); strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ctx.Query("token")
}

func NewDevAuthMiddleware() func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if isPublic(ctx.Operation()) {
			next(ctx)
			return
		}
		userID := ctx.Header("X-Dev-User-ID")
		if userID == "" {
			userID = "dev-user"
		}
		role := ctx.Header("X-Dev-Role")
		if role == "" {
			role = "admin"
		}
		email := ctx.Header("X-Dev-Email")
		if email == "" {
			email = "dev@example.test"
		}
		name := ctx.Header("X-Dev-Name")
		if name == "" {
			name = "Dev User"
		}
		claims := UserClaims{UserID: userID, Email: email, Role: role, Name: name}
		next(huma.WithContext(ctx, ContextWithClaims(ctx.Context(), claims)))
	}
}

func isPublic(op *huma.Operation) bool {
	if op == nil {
		return false
	}
	return op.Metadata != nil && op.Metadata["public"] == true
}

func PublicOperation(o *huma.Operation) {
	if o.Metadata == nil {
		o.Metadata = map[string]any{}
	}
	o.Metadata["public"] = true
}

func writeUnauthorized(ctx huma.Context, message string) {
	ctx.SetStatus(401)
	ctx.SetHeader("Content-Type", "application/json")
	_ = json.NewEncoder(ctx.BodyWriter()).Encode(map[string]string{"error": message})
}

func writeForbidden(ctx huma.Context, message string) {
	ctx.SetStatus(403)
	ctx.SetHeader("Content-Type", "application/json")
	_ = json.NewEncoder(ctx.BodyWriter()).Encode(map[string]string{"error": message})
}

// RequireRoles membatasi sebuah operasi Huma hanya untuk role tertentu.
// Dipakai bersama middleware autentikasi global yang sudah menaruh claims di
// context (lihat NewJWTAuthMiddleware / NewDevAuthMiddleware).
func RequireRoles(roles ...string) func(*huma.Operation) {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(o *huma.Operation) {
		o.Middlewares = append(o.Middlewares, func(ctx huma.Context, next func(huma.Context)) {
			claims, ok := ClaimsFromContext(ctx.Context())
			if !ok {
				writeUnauthorized(ctx, "unauthenticated")
				return
			}
			if _, ok := allowed[claims.Role]; !ok {
				writeForbidden(ctx, "forbidden")
				return
			}
			next(ctx)
		})
	}
}
