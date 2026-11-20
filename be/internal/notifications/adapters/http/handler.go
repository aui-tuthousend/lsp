package http

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2/sse"
	"lsp/internal/notifications/domain"
	"lsp/internal/platform/auth"
	"lsp/internal/platform/events"
)

// heartbeatInterval adalah jeda pengiriman komentar SSE untuk menjaga
// koneksi tetap hidup melewati proxy/load balancer.
const heartbeatInterval = 20 * time.Second

type Handler struct {
	bus      events.Bus
	shutdown <-chan struct{}
}

func NewHandler(bus events.Bus, shutdown <-chan struct{}) *Handler {
	return &Handler{bus: bus, shutdown: shutdown}
}

// Stream mengirim notifikasi realtime ke klien sampai koneksi ditutup atau
// aplikasi dihentikan.
func (h *Handler) Stream(ctx context.Context, input *EventsInput, send sse.Sender) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		_ = send.Data(domain.NotificationEvent{Type: "error", Message: "unauthenticated"})
		return
	}

	topics := []string{events.UserTopic(claims.UserID)}
	if claims.Role != "" {
		topics = append(topics, events.RoleTopic(claims.Role))
	}
	sub, err := h.bus.Subscribe(ctx, topics...)
	if err != nil {
		_ = send.Data(domain.NotificationEvent{Type: "error", Message: "gagal berlangganan event"})
		return
	}
	defer func() { _ = sub.Close() }()

	_ = send.Data(domain.NotificationEvent{Type: "connected", Message: "stream terbuka"})

	heartbeat := time.NewTicker(heartbeatInterval)
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-h.shutdown:
			return
		case <-heartbeat.C:
			if err := send.Comment("ping"); err != nil {
				return
			}
		case evt, ok := <-sub.C():
			if !ok {
				return
			}
			if err := send.Data(domain.NotificationEvent{Type: evt.Type, Topic: evt.Topic, Data: evt.Data}); err != nil {
				return
			}
		}
	}
}
