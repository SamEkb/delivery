package commands

import (
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidOrderId = errors.New("order id must not be Empty")
var ErrInvalidStreet = errors.New("street must not be empty")
var ErrInvalidVolume = errors.New("volume must be greater than zero")

type CreateOrderCommand struct {
	orderID uuid.UUID
	street  string
	volume  int

	isValid bool
}

func NewCreateOrderCommand(orderID uuid.UUID, street string, volume int) (*CreateOrderCommand, error) {
	if orderID == uuid.Nil {
		return nil, ErrInvalidOrderId
	}
	if street == "" {
		return nil, ErrInvalidStreet
	}
	if volume <= 0 {
		return nil, ErrInvalidVolume
	}
	return &CreateOrderCommand{
		orderID: orderID,
		street:  street,
		volume:  volume,
		isValid: true,
	}, nil
}

func (c *CreateOrderCommand) IsValid() bool {
	return c.isValid
}

func (c *CreateOrderCommand) OrderID() uuid.UUID {
	return c.orderID
}

func (c *CreateOrderCommand) Street() string {
	return c.street
}

func (c *CreateOrderCommand) Volume() int {
	return c.volume
}
