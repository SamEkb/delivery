package outbox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/delivery/internal/pkg/ddd"
	"github.com/delivery/internal/pkg/errs"
)

type EventRegistry interface {
	RegisterDomainEvent(eventType reflect.Type) error
	DecodeDomainEvent(event *Event) (ddd.DomainEvent, error)
}

var _ EventRegistry = &eventRegistry{}

type eventRegistry struct {
	events map[string]reflect.Type
}

func NewEventRegistry() (EventRegistry, error) {
	return &eventRegistry{
		events: make(map[string]reflect.Type),
	}, nil
}

func (r *eventRegistry) RegisterDomainEvent(eventType reflect.Type) error {
	if eventType == nil {
		return errs.NewValueIsRequiredError("event type")
	}
	eventName := eventType.Name()
	r.events[eventName] = eventType

	return nil
}

func EncodeDomainEvent(domainEvent ddd.DomainEvent) (Event, error) {
	payload, err := json.Marshal(domainEvent)
	if err != nil {
		return Event{}, fmt.Errorf("failed to encode event: %w", err)
	}

	return Event{
		ID:             domainEvent.GetID(),
		Name:           domainEvent.GetName(),
		Payload:        payload,
		OccurredAtUtc:  time.Now().UTC(),
		ProcessedAtUtc: nil,
	}, nil
}

func EncodeDomainEvents(domainEvents []ddd.DomainEvent) ([]Event, error) {
	if len(domainEvents) == 0 {
		return nil, errs.NewValueIsRequiredError("domain events")
	}

	events := make([]Event, 0, len(domainEvents))
	for _, domainEvent := range domainEvents {
		event, err := EncodeDomainEvent(domainEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to encode domain event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *eventRegistry) DecodeDomainEvent(event *Event) (ddd.DomainEvent, error) {
	if event == nil {
		return nil, errs.NewValueIsRequiredError("event")
	}

	eventType, ok := r.events[event.Name]
	if !ok {
		return nil, errs.NewNotFoundError("event", event.Name)
	}

	decodedEvent := reflect.New(eventType).Interface()

	if err := json.Unmarshal(event.Payload, decodedEvent); err != nil {
		return nil, fmt.Errorf("failed to decode event: %w", err)
	}

	domainEvent, ok := decodedEvent.(ddd.DomainEvent)
	if !ok {
		return nil, fmt.Errorf("invalid event type: %T, expected: %T", decodedEvent, ddd.DomainEvent(nil))
	}

	return domainEvent, nil
}
