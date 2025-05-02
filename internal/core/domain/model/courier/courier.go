package courier

import (
	"errors"
	"math"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
)

var ErrInvalidCourierName = errors.New("courier's name should not be empty")
var ErrInvalidSpeedValue = errors.New("speed must be greater than 0")
var ErrInvalidOrder = errors.New("order can't be null")
var ErrCanNotTakeOrder = errors.New("can't take order, all storage places are occupied")
var ErrOrderNotFound = errors.New("order not found in any storage place")

type Courier struct {
	id            uuid.UUID
	name          string
	speed         int
	location      kernel.Location
	storagePlaces []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, ErrInvalidCourierName
	}

	if speed <= 0 {
		return nil, ErrInvalidSpeedValue
	}

	return &Courier{
		id:            uuid.New(),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: make([]*StoragePlace, 0),
	}, nil
}

func (c *Courier) Equals(other *Courier) bool {
	if other == nil {
		return false
	}
	return c.id == other.id
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	storagePlace, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}

	c.storagePlaces = append(c.storagePlaces, storagePlace)

	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) (bool, error) {
	if order == nil {
		return false, ErrInvalidOrder
	}

	for _, v := range c.storagePlaces {
		canStore := v.CanStore(order.Volume())
		if canStore {
			return true, nil
		}
	}

	return false, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {
	if order == nil {
		return ErrInvalidOrder
	}

	canTake, err := c.CanTakeOrder(order)
	if err != nil {
		return err
	}

	if !canTake {
		return ErrCanNotTakeOrder
	}

	for _, v := range c.storagePlaces {
		canStore := v.CanStore(order.Volume())

		if canStore {
			err := v.Store(order.ID(), order.Volume())
			if err != nil {
				return err
			}
			return nil
		}

	}

	return ErrCanNotTakeOrder
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	if order == nil {
		return ErrInvalidOrder
	}

	if err := order.Complete(); err != nil {
		return err
	}

	orderID := order.ID()

	storagePlace, err := c.findStoragePlaceByOrderID(orderID)
	if err != nil {
		return err
	}

	if err := storagePlace.Clear(orderID); err != nil {
		return err
	}

	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) float64 {
	courierSpeed := c.speed
	courierLocation := c.location

	distance := target.DistanceTo(courierLocation)

	return float64(distance) / float64(courierSpeed)
}

func (c *Courier) Move(target kernel.Location) error {
	x := float64(target.X() - c.location.X())
	y := float64(target.Y() - c.location.Y())
	speed := float64(c.speed)

	if math.Abs(x) > speed {
		x = math.Copysign(speed, x)
	}
	speed -= math.Abs(x)

	if math.Abs(y) > speed {
		y = math.Copysign(speed, y)
	}

	newX := c.location.X() + int(x)
	newY := c.location.Y() + int(y)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation
	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	for _, v := range c.storagePlaces {
		if v.OrderID() != nil && *v.OrderID() == orderID {
			return v, nil
		}
	}

	return nil, ErrOrderNotFound
}

func (c *Courier) ID() uuid.UUID {
	return c.id
}
