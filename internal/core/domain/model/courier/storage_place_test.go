package courier

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewStoragePlace(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int
		expectedErr error
		wantErr     bool
	}{
		"valid storage place": {
			name:        "storage place",
			totalVolume: 10,
			expectedErr: nil,
			wantErr:     false,
		}, "invalid name": {
			name:        "",
			totalVolume: 10,
			expectedErr: ErrInvalidName,
			wantErr:     true,
		}, "invalid total volume": {
			name:        "storage place",
			totalVolume: 0,
			expectedErr: ErrInvalidTotalVolume,
			wantErr:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewStoragePlace(tc.name, tc.totalVolume)
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStoragePlace_CanStore(t *testing.T) {
	tests := map[string]struct {
		amount         int
		expectedResult bool
	}{
		"can store": {
			amount:         10,
			expectedResult: true,
		}, "can not store positive value": {
			amount:         11,
			expectedResult: false,
		}, "can not store negative value": {
			amount:         -1,
			expectedResult: false,
		}, "can not store zero": {
			amount:         0,
			expectedResult: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, _ := NewStoragePlace("storage place", 10)
			result := s.CanStore(tc.amount)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestStoragePlace_Store(t *testing.T) {
	tests := map[string]struct {
		orderID        uuid.UUID
		amount         int
		expectedResult error
		wantErr        bool
	}{
		"can store": {
			orderID:        uuid.New(),
			amount:         1,
			expectedResult: nil,
			wantErr:        false,
		},
		"invalid order id": {
			orderID:        uuid.Nil,
			amount:         1,
			expectedResult: ErrInvalidOrderId,
			wantErr:        true,
		}, "invalid amount": {
			orderID:        uuid.New(),
			amount:         0,
			expectedResult: ErrCanNotStore,
			wantErr:        true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, _ := NewStoragePlace("storage place", 10)
			err := s.Store(tc.orderID, tc.amount)
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.expectedResult)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStoragePlace_Clear(t *testing.T) {
	tests := map[string]struct {
		storedOrderID  uuid.UUID
		clearOrderID   uuid.UUID
		expectedResult error
		wantErr        bool
	}{
		"can clear": {
			storedOrderID:  uuid.New(),
			clearOrderID:   uuid.Nil,
			expectedResult: nil,
			wantErr:        false,
		}, "invalid order id": {
			storedOrderID:  uuid.New(),
			clearOrderID:   uuid.Nil,
			expectedResult: ErrInvalidOrderId,
			wantErr:        true,
		}, "wrong order id": {
			storedOrderID:  uuid.New(),
			clearOrderID:   uuid.New(),
			expectedResult: ErrWrongOrderId,
			wantErr:        true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, _ := NewStoragePlace("storage place", 10)
			s.orderID = &tc.storedOrderID

			clearID := tc.clearOrderID
			if name == "can clear" {
				clearID = tc.storedOrderID
			}

			err := s.Clear(clearID)
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.expectedResult)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStoragePlace_Equals(t *testing.T) {
	sharedID := uuid.New()

	tests := map[string]struct {
		setupTestObject func() (*StoragePlace, *StoragePlace)
		expectedResult  bool
	}{
		"equals": {
			setupTestObject: func() (*StoragePlace, *StoragePlace) {
				s1, _ := NewStoragePlace("storage place", 10)
				s1.id = sharedID

				s2 := &StoragePlace{
					id:          sharedID,
					name:        "storage place",
					totalVolume: 20,
				}
				return s1, s2
			},
			expectedResult: true,
		},
		"not equals": {
			setupTestObject: func() (*StoragePlace, *StoragePlace) {
				s1, _ := NewStoragePlace("storage place", 10)

				s2 := &StoragePlace{
					id:          uuid.New(),
					name:        "storage place",
					totalVolume: 10,
				}
				return s1, s2
			},
			expectedResult: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s1, s2 := tc.setupTestObject()
			result := s1.Equals(s2)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
