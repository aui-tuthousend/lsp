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
