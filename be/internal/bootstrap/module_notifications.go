package bootstrap

import (
	"context"

	notificationshttp "lsp/internal/features/notifications/adapters/http"
	"lsp/internal/platform/events"
)

type NotificationsModule struct {
	Bus     events.Bus
	Handler *notificationshttp.Handler
}

func newNotificationsModule(ctx context.Context, bus events.Bus) *NotificationsModule {
	return &NotificationsModule{Bus: bus, Handler: notificationshttp.NewHandler(bus, ctx.Done())}
}
