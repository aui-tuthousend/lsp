package http

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"lsp/internal/notifications/domain"
	httpplatform "lsp/internal/platform/http"
)

// eventTypeNotification adalah nama event SSE untuk NotificationEvent.
const eventTypeNotification = "notification"

func RegisterRoutes(api huma.API, handler *Handler) {
	op := huma.Operation{
		Method: "GET",
		Path:   "/events",
		Errors: []int{401},
	}
	httpplatform.WithDoc(
		"events-stream",
		"Subscribe realtime events",
		"Stream Server-Sent Events berisi notifikasi realtime untuk user yang sedang login. Autentikasi via header `Authorization: Bearer` atau query `?token=`.",
		"Notifications",
	)(&op)

	sse.Register(api, op, map[string]any{
		eventTypeNotification: domain.NotificationEvent{},
	}, handler.Stream)
}
