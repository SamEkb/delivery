package ddd

import (
	"context"
	"sync"
)

type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

type Mediatr interface {
	Subscribe(handler EventHandler, events ...DomainEvent)
	Publish(ctx context.Context, event DomainEvent) error
}

type mediatr struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

func NewMediatr() Mediatr {
	return &mediatr{handlers: make(map[string][]EventHandler)}
}

func (m *mediatr) Subscribe(handler EventHandler, events ...DomainEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, event := range events {
		handlers := m.handlers[event.GetName()]
		handlers = append(handlers, handler)
		m.handlers[event.GetName()] = handlers
	}
}

func (m *mediatr) Publish(ctx context.Context, event DomainEvent) error {

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, handler := range m.handlers[event.GetName()] {
		err := handler.Handle(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}
