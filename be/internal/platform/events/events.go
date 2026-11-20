// Package events menyediakan abstraksi pub/sub lintas modul untuk fan-out
// event real-time (mis. Server-Sent Events). Implementasi produksi memakai
// Redis Pub/Sub agar bisa dipakai lintas instance API.
package events

import (
	"context"
	"encoding/json"
)

// Event adalah satu pesan yang dipublikasikan ke sebuah topic.
type Event struct {
	Topic string          // mis. "user:<id>" atau "role:admin"
	Type  string          // mis. "apl01.approved"
	Data  json.RawMessage // payload JSON
}

// Subscription adalah hasil berlangganan sejumlah topic.
type Subscription interface {
	// C mengembalikan channel event. Channel ditutup saat Close dipanggil.
	C() <-chan Event
	// Close menghentikan langganan dan melepas resource.
	Close() error
}

// Bus adalah port pub/sub. Publisher dan subscriber tidak saling mengenal.
type Bus interface {
	Publish(ctx context.Context, evt Event) error
	Subscribe(ctx context.Context, topics ...string) (Subscription, error)
}

// UserTopic mengembalikan nama topic untuk notifikasi personal seorang user.
func UserTopic(userID string) string { return "user:" + userID }

// RoleTopic mengembalikan nama topic untuk notifikasi berbasis role.
func RoleTopic(role string) string { return "role:" + role }
