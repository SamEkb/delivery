package courier

import (
	"math"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/pkg/ddd"
	"github.com/delivery/internal/pkg/errs"
	"github.com/google/uuid"
)

type Courier struct {
	*ddd.BaseAggregate[uuid.UUID]
	name          string
	speed         int
	location      kernel.Location
	storagePlaces []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("courier name")
	}

	if speed <= 0 {
		return nil, errs.NewValueIsRequiredError("courier speed")
	}

	courierID := uuid.New()
	return &Courier{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](courierID),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: make([]*StoragePlace, 0),
	}, nil
}

func RestoreCourier(id uuid.UUID, name string, speed int, location kernel.Location, storagePlaces []*StoragePlace) *Courier {
	return &Courier{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](id),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: storagePlaces,
	}
}

func (c *Courier) Equals(other *Courier) bool {
	if other == nil {
		return false
	}
	return c.BaseAggregate.ID() == other.BaseAggregate.ID()
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
		return false, errs.NewValueIsRequiredError("order")
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
		return errs.NewValueIsRequiredError("order")
	}

	canTake, err := c.CanTakeOrder(order)
	if err != nil {
		return err
	}

	if !canTake {
		return errs.NewBusinessError("courier can't take order", "can't take order")
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

	return errs.NewBusinessError("courier can't take order", "all storage places are occupied")
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequiredError("order")
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

	return nil, errs.NewBusinessError("courier can't take order", "order is not stored here")
}

func (c *Courier) ID() uuid.UUID {
	return c.BaseAggregate.ID()
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) StoragePlaces() []*StoragePlace {
	return c.storagePlaces
}
