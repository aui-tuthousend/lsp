// Package domain berisi tipe data notifikasi realtime yang dikirim ke klien
// melalui Server-Sent Events.
package domain

import "encoding/json"

// NotificationEvent adalah payload SSE. Nama event pada stream selalu
// "notification"; jenis sebenarnya ada pada field Type.
type NotificationEvent struct {
	Type    string          `json:"type"`
	Topic   string          `json:"topic,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}
