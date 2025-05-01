package order

import (
	"testing"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	tests := map[string]struct {
		orderID  uuid.UUID
		location kernel.Location
		volume   int
		wantErr  bool
		err      error
	}{
		"valid order": {
			orderID:  uuid.New(),
			location: mustCreateLocation(1, 1),
			volume:   1,
			wantErr:  false,
		},
		"invalid orderID": {
			orderID:  uuid.Nil,
			location: mustCreateLocation(1, 1),
			volume:   1,
			wantErr:  true,
			err:      ErrInvalidOrderId,
		},
		"invalid order volume 0": {
			orderID:  uuid.New(),
			location: mustCreateLocation(1, 1),
			volume:   0,
			wantErr:  true,
			err:      ErrInvalidOrderVolume,
		},
		"invalid order volume negative": {
			orderID:  uuid.New(),
			location: mustCreateLocation(1, 1),
			volume:   -1,
			wantErr:  true,
			err:      ErrInvalidOrderVolume,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			order, err := NewOrder(tc.orderID, tc.location, tc.volume)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.orderID, order.ID())
				assert.Equal(t, tc.location, order.Location())
				assert.Equal(t, tc.volume, order.Volume())
				assert.Equal(t, created, order.Status())
				assert.Nil(t, order.CourierID())
			}
		})
	}
}

func TestOrder_Assign(t *testing.T) {
	validCourierID := uuid.New()
	tests := map[string]struct {
		courierId *uuid.UUID
		status    Status
		wantErr   bool
		err       error
	}{
		"valid assignation": {
			courierId: &validCourierID,
			status:    assigned,
			wantErr:   false,
			err:       nil,
		},
		"invalid courierID": {
			courierId: nil,
			status:    created,
			wantErr:   true,
			err:       ErrInvalidCourierId,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			validOrderID := uuid.New()
			location := mustCreateLocation(1, 1)
			order, err := NewOrder(validOrderID, location, 1)
			assert.NoError(t, err)

			err = order.Assign(tc.courierId)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
				assert.Equal(t, tc.status, order.Status())
				assert.Equal(t, tc.courierId, order.CourierID())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.status, order.Status())
				assert.Equal(t, tc.courierId, order.CourierID())
			}
		})
	}
}

func TestOrder_Complete(t *testing.T) {
	validCourierId := uuid.New()
	tests := map[string]struct {
		courierId *uuid.UUID
		status    Status
		wantErr   bool
		err       error
	}{
		"valid completion": {
			courierId: &validCourierId,
			status:    completed,
			wantErr:   false,
			err:       nil,
		},
		"courier didn't assign": {
			courierId: nil,
			status:    created,
			wantErr:   true,
			err:       ErrCourierWasNotAssign,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			validOrderID := uuid.New()
			location := mustCreateLocation(1, 1)
			order, err := NewOrder(validOrderID, location, 1)
			assert.NoError(t, err)

			if !tc.wantErr {
				err = order.Assign(tc.courierId)
			}

			err = order.Complete()

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
				assert.Equal(t, tc.status, order.Status())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.courierId, order.CourierID())
				assert.Equal(t, tc.status, order.Status())
			}
		})
	}
}

func TestOrder_Equals(t *testing.T) {
	validOrderID := uuid.New()
	tests := map[string]struct {
		order    *Order
		other    *Order
		expected bool
	}{
		"orders are equal": {
			order:    mustCreateOrder(validOrderID),
			other:    mustCreateOrder(validOrderID),
			expected: true,
		},
		"orders are different": {
			order:    mustCreateOrder(validOrderID),
			other:    mustCreateOrder(uuid.New()),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.order.Equals(tc.other)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func mustCreateOrder(orderID uuid.UUID) *Order {
	location := mustCreateLocation(1, 1)
	order, err := NewOrder(orderID, location, 1)
	if err != nil {
		panic(err)
	}

	return order
}

func mustCreateLocation(x, y int) kernel.Location {
	loc, err := kernel.NewLocation(x, y)
	if err != nil {
		panic(err)
	}
	return loc
}
