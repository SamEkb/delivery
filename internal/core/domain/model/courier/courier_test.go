package courier

import (
	"testing"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewCourier(t *testing.T) {
	validLocation := mustCreateLocation(1, 1)
	tests := map[string]struct {
		name     string
		speed    int
		location kernel.Location
		wantErr  bool
		err      error
	}{
		"valid courier": {
			name:     "courier12",
			speed:    2,
			location: validLocation,
			wantErr:  false,
			err:      nil,
		},
		"invalid couriers name": {
			name:     "",
			speed:    2,
			location: validLocation,
			wantErr:  true,
			err:      ErrInvalidCourierName,
		},
		"invalid couriers speed": {
			name:     "courier12",
			speed:    0,
			location: validLocation,
			wantErr:  true,
			err:      ErrInvalidSpeedValue,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			courier, err := NewCourier(tc.name, tc.speed, tc.location)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, courier.name, tc.name)
				assert.Equal(t, courier.speed, tc.speed)
				assert.Equal(t, courier.location, tc.location)
			}
		})
	}
}

func TestCourier_AddStoragePlace(t *testing.T) {
	validLocation := mustCreateLocation(1, 1)
	courier, err := NewCourier("courier12", 2, validLocation)
	assert.NoError(t, err)
	assert.NotNil(t, courier)
	err = courier.AddStoragePlace("storage place", 10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(courier.storagePlaces), 1)
}

func TestCourier_CanTakeOrder(t *testing.T) {
	validOrder := mustCreateOrder(uuid.New())
	tests := map[string]struct {
		order   *order.Order
		result  bool
		wantErr bool
		err     error
	}{
		"can take order": {
			order:   validOrder,
			result:  true,
			wantErr: false,
			err:     nil,
		},
		"cant take, no storage place": {
			order:   validOrder,
			result:  false,
			wantErr: false,
			err:     nil,
		},
		"cant take, invalid order": {
			order:   nil,
			result:  false,
			wantErr: true,
			err:     ErrInvalidOrder,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			courier, err := NewCourier("courier12", 2, mustCreateLocation(1, 1))
			assert.NoError(t, err)
			if tc.result {
				err = courier.AddStoragePlace("storage place", 10)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(courier.storagePlaces), 1)
			}
			canTake, err := courier.CanTakeOrder(tc.order)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.result, canTake)
			}
		})
	}
}

func TestCourier_TakeOrder(t *testing.T) {
	validOrder := mustCreateOrder(uuid.New())
	tests := map[string]struct {
		order    *order.Order
		canStore bool
		wantErr  bool
		err      error
	}{
		"can take order": {
			order:    validOrder,
			canStore: true,
			wantErr:  false,
			err:      nil,
		},
		"cant take, no storage place": {
			order:    validOrder,
			canStore: false,
			wantErr:  true,
			err:      ErrCanNotTakeOrder,
		},
		"cant take, invalid order": {
			order:    nil,
			canStore: false,
			wantErr:  true,
			err:      ErrInvalidOrder,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			courier, err := NewCourier("courier12", 2, mustCreateLocation(1, 1))
			assert.NoError(t, err)
			if tc.canStore {
				err = courier.AddStoragePlace("storage place", 10)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(courier.storagePlaces), 1)
			}
			err = courier.TakeOrder(tc.order)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCourier_CompleteOrder(t *testing.T) {
	ord := mustCreateOrder(uuid.New())
	tests := map[string]struct {
		order   *order.Order
		wantErr bool
		err     error
	}{
		"can successfully complete order": {
			order:   ord,
			wantErr: false,
			err:     nil,
		},
		"cant complete order, invalid order": {
			order:   nil,
			wantErr: true,
			err:     ErrInvalidOrder,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			courier, err := NewCourier("courier12", 2, mustCreateLocation(1, 1))
			assert.NoError(t, err)

			if !tc.wantErr {
				courierID := courier.ID()
				err = ord.Assign(&courierID)
				assert.NoError(t, err)

				err = courier.AddStoragePlace("storage place", 10)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(courier.storagePlaces), 1)

				err = courier.TakeOrder(tc.order)
				assert.NoError(t, err)
			}

			err = courier.CompleteOrder(tc.order)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, ord.Status(), order.Completed)
			}
		})
	}
}

func TestCourier_CalculateTimeToLocation(t *testing.T) {
	startLocation := mustCreateLocation(1, 1)
	courier, err := NewCourier("courier12", 2, startLocation)
	assert.NoError(t, err)
	assert.NotNil(t, courier)

	endLocation := mustCreateLocation(3, 3)

	time := courier.CalculateTimeToLocation(endLocation)
	assert.NotNil(t, time)
}

func TestCourier_Move(t *testing.T) {
	tests := map[string]struct {
		initialX     int
		initialY     int
		courierSpeed int
		targetX      int
		targetY      int
		expectedX    int
		expectedY    int
		wantErr      bool
	}{
		"move within speed limit": {
			initialX:     3,
			initialY:     3,
			courierSpeed: 5,
			targetX:      6,
			targetY:      5,
			expectedX:    6,
			expectedY:    5,
			wantErr:      false,
		},
		"move exceeding speed limit on both axes": {
			initialX:     2,
			initialY:     2,
			courierSpeed: 3,
			targetX:      10,
			targetY:      10,
			expectedX:    5,
			expectedY:    2,
			wantErr:      false,
		},
		"move exceeding speed limit on Y axis": {
			initialX:     2,
			initialY:     2,
			courierSpeed: 5,
			targetX:      4,
			targetY:      10,
			expectedX:    4,
			expectedY:    5,
			wantErr:      false,
		},
		"no movement needed": {
			initialX:     5,
			initialY:     5,
			courierSpeed: 10,
			targetX:      5,
			targetY:      5,
			expectedX:    5,
			expectedY:    5,
			wantErr:      false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			initialLocation, err := kernel.NewLocation(tc.initialX, tc.initialY)
			assert.NoError(t, err)

			courier, err := NewCourier("test-courier", tc.courierSpeed, initialLocation)
			assert.NoError(t, err)

			targetLocation, err := kernel.NewLocation(tc.targetX, tc.targetY)
			assert.NoError(t, err)

			err = courier.Move(targetLocation)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedX, courier.location.X())
				assert.Equal(t, tc.expectedY, courier.location.Y())
			}
		})
	}
}

func mustCreateOrder(orderID uuid.UUID) *order.Order {
	location := mustCreateLocation(1, 1)
	ord, err := order.NewOrder(orderID, location, 1)
	if err != nil {
		panic(err)
	}
	return ord
}

func mustCreateLocation(x, y int) kernel.Location {
	loc, err := kernel.NewLocation(x, y)
	if err != nil {
		panic(err)
	}
	return loc
}
