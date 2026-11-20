package http

import (
	"context"
	"log/slog"
	stdhttp "net/http"
	"runtime/debug"
	"time"

	"github.com/rs/cors"
	"lsp/internal/platform/id"
)

type Middleware func(stdhttp.Handler) stdhttp.Handler

type requestIDKey struct{}

func Chain(handler stdhttp.Handler, middlewares ...Middleware) stdhttp.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func RequestID(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = id.New()
		}
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey{}).(string)
	return value
}

func Logger(logger *slog.Logger) Middleware {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.InfoContext(r.Context(), "http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String(), "request_id", RequestIDFromContext(r.Context()))
		})
	}
}

func Recover(logger *slog.Logger) Middleware {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.ErrorContext(r.Context(), "panic recovered", "error", rec, "stack", string(debug.Stack()))
					stdhttp.Error(w, "internal server error", stdhttp.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func CORS() Middleware {
	c := cors.New(cors.Options{AllowedOrigins: []string{"*"}, AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, AllowedHeaders: []string{"Authorization", "Content-Type", "X-Request-ID", "X-Dev-User-ID", "X-Dev-Role", "X-Dev-Email", "X-Dev-Name"}})
	return c.Handler
}

func MaxBodySize(bytes int64) Middleware {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			r.Body = stdhttp.MaxBytesReader(w, r.Body, bytes)
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout membatasi durasi request. Path pada skipPaths dilewatkan tanpa
// timeout karena responsnya bersifat streaming (mis. SSE) dan tidak boleh
// dibuffer/dipotong oleh http.TimeoutHandler.
func Timeout(timeout time.Duration, skipPaths ...string) Middleware {
	skip := make(map[string]struct{}, len(skipPaths))
	for _, path := range skipPaths {
		skip[path] = struct{}{}
	}
	return func(next stdhttp.Handler) stdhttp.Handler {
		limited := stdhttp.TimeoutHandler(next, timeout, "request timeout")
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			if _, ok := skip[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}
			limited.ServeHTTP(w, r)
		})
	}
}
