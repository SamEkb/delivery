package order

import (
	"reflect"

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

func NewStatusChangedDomainEvent(order *Order) (*StatusChangedDomainEvent, error) {
	event := &StatusChangedDomainEvent{
		ID: uuid.New(),

		OrderID:     order.ID(),
		OrderStatus: order.Status(),

		isValid: true,
	}
	event.Name = reflect.TypeOf(event).Name()

	return event, nil
}

func NewStatusChangedDomainEventWithoutData() *StatusChangedDomainEvent {
	return &StatusChangedDomainEvent{}
}
