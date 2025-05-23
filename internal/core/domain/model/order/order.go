package order

import (
	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/pkg/ddd"
	"github.com/delivery/internal/pkg/errs"
	"github.com/google/uuid"
)

type Order struct {
	*ddd.BaseAggregate[uuid.UUID]

	courierID *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("order id")
	}

	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}

	return &Order{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](orderID),
		courierID:     nil,
		location:      location,
		volume:        volume,
		status:        Created,
	}, nil
}

// RestoreOrder must be used ONLY in a repository layer for mapping
func RestoreOrder(orderID uuid.UUID, courierID *uuid.UUID, location kernel.Location, volume int, status Status) *Order {
	return &Order{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](orderID),
		courierID:     courierID,
		location:      location,
		volume:        volume,
		status:        status,
	}
}

func (o *Order) Assign(courierId *uuid.UUID) error {
	if courierId == nil {
		return errs.NewValueIsRequiredError("courier id")
	}

	o.courierID = courierId
	o.status = Assigned

	return nil
}

func (o *Order) Complete() error {
	if o.courierID == nil {
		return errs.NewBusinessError("order is not assigned to courier", "courier id is nil")
	}

	o.status = Completed

	domainEvent, err := NewStatusChangedDomainEvent(o.ID(), "", o)
	if err != nil {
		return err
	}
	o.BaseAggregate.RaiseDomainEvent(domainEvent)

	return nil
}

func (o *Order) Equals(other *Order) bool {
	return o.BaseAggregate.ID() == other.BaseAggregate.ID()
}

func (o *Order) ID() uuid.UUID {
	return o.BaseAggregate.ID()
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}
