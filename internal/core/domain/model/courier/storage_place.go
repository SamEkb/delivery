package courier

import (
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidName = errors.New("name is empty")
var ErrInvalidTotalVolume = errors.New("total volume must be greater than 0")
var ErrInvalidOrderId = errors.New("order id is empty")
var ErrCanNotStore = errors.New("order can't be stored")
var ErrWrongOrderId = errors.New("this order is not stored here")

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	if totalVolume <= 0 {
		return nil, ErrInvalidTotalVolume
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
		orderID:     &uuid.Nil,
	}, nil
}

func (s *StoragePlace) CanStore(amount int) bool {
	if amount <= 0 {
		return false
	}
	return s.totalVolume >= amount && !s.isOccupied()
}

func (s *StoragePlace) Store(orderID uuid.UUID, amount int) error {
	if orderID == uuid.Nil {
		return ErrInvalidOrderId
	}
	if !s.CanStore(amount) {
		return ErrCanNotStore
	}

	s.orderID = &orderID

	return nil
}

func (s *StoragePlace) Clear(orderID uuid.UUID) error {
	if orderID == uuid.Nil {
		return ErrInvalidOrderId
	}

	if s.orderID == nil || *s.orderID != orderID {
		return ErrWrongOrderId
	}

	*s.orderID = uuid.Nil

	return nil
}

func (s *StoragePlace) isOccupied() bool {
	return *s.orderID != uuid.Nil
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
	return s.id == other.id
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}
