package order

import (
	"github.com/delivery/internal/pkg/ddd"
	"github.com/google/uuid"
)

var _ ddd.DomainEvent = &StatusChangedDomainEvent{}

type StatusChangedDomainEvent struct {
	// base
	ID   uuid.UUID
	Name string

	// payload
	OrderID     uuid.UUID
	OrderStatus Status

	isValid bool
}

func (o *StatusChangedDomainEvent) GetID() uuid.UUID {
	return o.ID
}

func (o *StatusChangedDomainEvent) GetName() string {
	return o.Name
}

func NewStatusChangedDomainEvent(id uuid.UUID, name string, payload *Order) (*StatusChangedDomainEvent, error) {
	return &StatusChangedDomainEvent{
		ID:          id,
		Name:        name,
		OrderID:     payload.ID(),
		OrderStatus: payload.Status(),
		isValid:     true,
	}, nil
}

func NewStatusChangedDomainEventWithoutData() *StatusChangedDomainEvent {
	return &StatusChangedDomainEvent{}
}
