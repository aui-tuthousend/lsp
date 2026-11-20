package events

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/redis/go-redis/v9"
)

// wireMessage adalah format pesan yang dikirim melalui Redis Pub/Sub.
type wireMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// RedisBus adalah adapter Bus berbasis Redis Pub/Sub.
type RedisBus struct {
	rdb *redis.Client
}

func NewRedisBus(rdb *redis.Client) *RedisBus {
	return &RedisBus{rdb: rdb}
}

func (b *RedisBus) Publish(ctx context.Context, evt Event) error {
	payload, err := json.Marshal(wireMessage{Type: evt.Type, Data: evt.Data})
	if err != nil {
		return err
	}
	return b.rdb.Publish(ctx, evt.Topic, payload).Err()
}

func (b *RedisBus) Subscribe(ctx context.Context, topics ...string) (Subscription, error) {
	if len(topics) == 0 {
		return nil, errors.New("events: minimal satu topic diperlukan")
	}
	ps := b.rdb.Subscribe(ctx, topics...)
	// Pastikan langganan benar-benar aktif sebelum dikembalikan, sehingga
	// event yang dipublikasikan setelah Subscribe tidak hilang.
	if _, err := ps.Receive(ctx); err != nil {
		_ = ps.Close()
		return nil, err
	}
	sub := &redisSubscription{ch: make(chan Event, 64), ps: ps, done: make(chan struct{})}
	go sub.run()
	return sub, nil
}

type redisSubscription struct {
	ch   chan Event
	ps   *redis.PubSub
	done chan struct{}
	once sync.Once
}

func (s *redisSubscription) C() <-chan Event { return s.ch }

func (s *redisSubscription) Close() error {
	s.once.Do(func() { close(s.done) })
	return s.ps.Close()
}

func (s *redisSubscription) run() {
	defer close(s.ch)
	for msg := range s.ps.Channel() {
		var wm wireMessage
		if err := json.Unmarshal([]byte(msg.Payload), &wm); err != nil {
			continue
		}
		select {
		case s.ch <- Event{Topic: msg.Channel, Type: wm.Type, Data: wm.Data}:
		case <-s.done:
			return
		}
	}
}
