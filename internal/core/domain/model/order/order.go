package order

import (
	"errors"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/google/uuid"
)

var ErrInvalidOrderId = errors.New("order id must not be empty")
var ErrInvalidOrderVolume = errors.New("volume order should be greater than zero")
var ErrInvalidCourierId = errors.New("courier id must not be empty")
var ErrCourierWasNotAssign = errors.New("can't complete the order, courier was not assigned")

type Order struct {
	id        uuid.UUID
	courierID *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, ErrInvalidOrderId
	}

	if volume <= 0 {
		return nil, ErrInvalidOrderVolume
	}

	return &Order{
		id:        orderID,
		courierID: nil,
		location:  location,
		volume:    volume,
		status:    created,
	}, nil
}

func (o *Order) Assign(courierId *uuid.UUID) error {
	if courierId == nil {
		return ErrInvalidCourierId
	}

	o.courierID = courierId
	o.status = assigned

	return nil
}

func (o *Order) Complete() error {
	if o.courierID == nil {
		return ErrCourierWasNotAssign
	}

	o.status = completed

	return nil
}

func (o *Order) Equals(other *Order) bool {
	return o.id == other.id
}

func (o *Order) ID() uuid.UUID {
	return o.id
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
