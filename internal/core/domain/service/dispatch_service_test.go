package service

import (
	"testing"

	"github.com/delivery/internal/core/domain/model/courier"
	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDispatchService_Dispatch(t *testing.T) {
	couriers := createCouriers()
	courier1 := couriers[0]
	courier2 := couriers[1]

	validOrder := mustCreateOrder(uuid.New())
	invalidOrder := mustCreateOrder(uuid.New())
	courierID := courier1.ID()
	err := invalidOrder.Assign(&courierID)
	assert.NoError(t, err)

	tests := map[string]struct {
		orderParam *order.Order
		couriers   []*courier.Courier
		result     *courier.Courier
		wantErr    bool
		err        error
	}{
		"valid order dispatching": {
			orderParam: validOrder,
			couriers:   couriers,
			result:     courier2,
			wantErr:    false,
			err:        nil,
		},
		"invalid order nil": {
			orderParam: nil,
			couriers:   couriers,
			result:     nil,
			wantErr:    true,
			err:        ErrInvalidOrder,
		},
		"invalid order wrong status": {
			orderParam: invalidOrder,
			couriers:   couriers,
			result:     nil,
			wantErr:    true,
			err:        ErrInvalidOrder,
		},
		"invalid couriers nil": {
			orderParam: validOrder,
			couriers:   nil,
			result:     nil,
			wantErr:    true,
			err:        ErrInvalidCouriers,
		},
		"no suitable couriers": {
			orderParam: validOrder,
			couriers:   createOccupiedCouriers(),
			result:     nil,
			wantErr:    true,
			err:        ErrCourierNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			service := NewDispatchService()
			dispatch, err := service.Dispatch(tc.orderParam, tc.couriers)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.result, dispatch)
			}
		})
	}
}

func mustCreateOrder(orderID uuid.UUID) *order.Order {
	location := mustCreateLocation(4, 4)
	ord, err := order.NewOrder(orderID, location, 1)
	if err != nil {
		panic(err)
	}
	return ord
}

func createOccupiedCouriers() []*courier.Courier {
	courier1 := mustCreateCourier("courier1", 1, mustCreateLocation(1, 1))
	mustAddStoragePlace(courier1, "storage place 1", 1)

	order1 := mustCreateOrder(uuid.New())
	courier1ID := courier1.ID()
	err := order1.Assign(&courier1ID)
	if err != nil {
		panic(err)
	}

	err = courier1.TakeOrder(order1)
	if err != nil {
		panic(err)
	}

	return []*courier.Courier{courier1}
}

func createCouriers() []*courier.Courier {
	courier1 := mustCreateCourier("courier1", 1, mustCreateLocation(1, 1))
	mustAddStoragePlace(courier1, "storage place 1", 10)
	courier2 := mustCreateCourier("courier2", 2, mustCreateLocation(2, 2))
	mustAddStoragePlace(courier2, "storage place 2", 10)

	return []*courier.Courier{courier1, courier2}
}

func mustCreateCourier(name string, speed int, location kernel.Location) *courier.Courier {
	newCourier, err := courier.NewCourier(name, speed, location)
	if err != nil {
		panic(err)
	}
	return newCourier
}

func mustAddStoragePlace(courier *courier.Courier, name string, capacity int) {
	err := courier.AddStoragePlace(name, capacity)
	if err != nil {
		panic(err)
	}
}

func mustCreateLocation(x, y int) kernel.Location {
	loc, err := kernel.NewLocation(x, y)
	if err != nil {
		panic(err)
	}
	return loc
}
