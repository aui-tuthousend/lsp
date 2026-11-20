package http

// EventsInput tidak butuh body. Token autentikasi diambil oleh middleware
// dari header Authorization atau query `?token=` (karena EventSource browser
// tidak bisa mengirim header custom).
type EventsInput struct{}
