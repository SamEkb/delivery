package kernel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocation(t *testing.T) {
	tests := map[string]struct {
		x, y     int
		expected Location
		wantErr  bool
	}{
		"valid location": {
			x:        5,
			y:        5,
			expected: Location{x: 5, y: 5},
			wantErr:  false,
		},
		"x too small": {
			x:       0,
			y:       5,
			wantErr: true,
		},
		"y too small": {
			x:       5,
			y:       0,
			wantErr: true,
		},
		"x too large": {
			x:       11,
			y:       5,
			wantErr: true,
		},
		"y too large": {
			x:       5,
			y:       11,
			wantErr: true,
		},
		"boundary valid x,y min": {
			x:        1,
			y:        1,
			expected: Location{x: 1, y: 1},
			wantErr:  false,
		},
		"boundary valid x,y max": {
			x:        10,
			y:        10,
			expected: Location{x: 10, y: 10},
			wantErr:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			loc, err := NewLocation(tc.x, tc.y)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidLocation)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, loc)
			}
		})
	}
}

func TestCreateRandomLocation(t *testing.T) {
	for i := 0; i < 10; i++ {
		loc := CreateRandomLocation()
		assert.GreaterOrEqual(t, loc.X(), 1)
		assert.LessOrEqual(t, loc.X(), 10)
		assert.GreaterOrEqual(t, loc.Y(), 1)
		assert.LessOrEqual(t, loc.Y(), 10)
	}
}

func TestLocation_DistanceTo(t *testing.T) {
	tests := map[string]struct {
		loc      Location
		other    Location
		expected int
	}{
		"same location": {
			loc:      Location{x: 5, y: 5},
			other:    Location{x: 5, y: 5},
			expected: 0,
		},
		"horizontal distance": {
			loc:      Location{x: 1, y: 5},
			other:    Location{x: 5, y: 5},
			expected: 4,
		},
		"vertical distance": {
			loc:      Location{x: 5, y: 1},
			other:    Location{x: 5, y: 5},
			expected: 4,
		},
		"diagonal distance": {
			loc:      Location{x: 1, y: 1},
			other:    Location{x: 5, y: 5},
			expected: 8,
		},
		"negative diff x": {
			loc:      Location{x: 5, y: 5},
			other:    Location{x: 1, y: 5},
			expected: 4,
		},
		"negative diff y": {
			loc:      Location{x: 5, y: 5},
			other:    Location{x: 5, y: 1},
			expected: 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dist := tc.loc.DistanceTo(tc.other)
			assert.Equal(t, tc.expected, dist)
		})
	}
}

func TestLocation_Equals(t *testing.T) {
	tests := map[string]struct {
		loc      Location
		other    Location
		expected bool
	}{
		"same location": {
			loc:      Location{x: 5, y: 5},
			other:    Location{x: 5, y: 5},
			expected: true,
		},
		"different x": {
			loc:      Location{x: 1, y: 5},
			other:    Location{x: 5, y: 5},
			expected: false,
		},
		"different y": {
			loc:      Location{x: 5, y: 1},
			other:    Location{x: 5, y: 5},
			expected: false,
		},
		"different x and y": {
			loc:      Location{x: 1, y: 1},
			other:    Location{x: 5, y: 5},
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			isEquals := tc.loc.Equals(tc.other)
			assert.Equal(t, tc.expected, isEquals)
		})
	}
}

func TestLocation_X_Y(t *testing.T) {
	loc := Location{x: 5, y: 10}

	assert.Equal(t, 5, loc.X())
	assert.Equal(t, 10, loc.Y())
}
