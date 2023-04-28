package pubsub

import (
	"context"
	"github.com/anaregdesign/papaya/model"
	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

type Subscription[T any] struct {
	mu          sync.RWMutex
	wg          sync.WaitGroup
	name        string
	topic       *Topic[T]
	ch          chan string
	messages    map[string]*Message[T]
	concurrency int64
	interval    time.Duration
	ttl         time.Duration
}

func (s *Subscription[T]) Name() string {
	return s.name
}

func (s *Subscription[T]) Topic() *Topic[T] {
	return s.topic
}

func (s *Subscription[T]) Subscribe(ctx context.Context, consumer model.Consumer[*Message[T]]) {
	sem := semaphore.NewWeighted(s.concurrency)
	s.register()
	go s.watch(ctx, s.interval, s.ttl)

	for {
		select {
		case id := <-s.ch:
			message := s.messages[id]
			s.wg.Add(1)
			sem.Acquire(ctx, 1)
			go func(m *Message[T]) {
				sem.Release(1)
				s.wg.Done()
				consumer(message)
			}(message)

		case <-ctx.Done():
			s.wg.Wait()
			s.unregister()
		}
	}
}

func (s *Subscription[T]) register() {
	s.topic.register(s)
}

func (s *Subscription[T]) unregister() {
	s.topic.unregister(s)
}

func (s *Subscription[T]) newMessage(body T) *Message[T] {
	return &Message[T]{
		id:           uuid.New().String(),
		body:         body,
		subscription: s,
		createdAt:    time.Now(),
	}
}

func (s *Subscription[T]) ack(message *Message[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.messages, message.id)
}

func (s *Subscription[T]) nack(message *Message[T]) {
}

func (s *Subscription[T]) remind(message *Message[T]) {
	s.ch <- message.id
}

func (s *Subscription[T]) salvage(interval time.Duration, ttl time.Duration) {
	for _, message := range s.messages {
		if time.Now().Sub(message.createdAt) > ttl {
			s.ack(message)
		}

		if time.Now().Sub(message.lastViewedAt) > interval {
			s.remind(message)
		}
	}
}

func (s *Subscription[T]) watch(ctx context.Context, interval time.Duration, ttl time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			s.salvage(interval, ttl)
		case <-ctx.Done():
			s.wg.Wait()
			return
		}
	}
}
