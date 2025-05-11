package kernel

import (
	"errors"
	"math"
	"math/rand"
)

const (
	minX = 1
	maxX = 10

	minY = 1
	maxY = 10
)

var ErrInvalidLocation = errors.New("location is invalid")

type Location struct {
	x int
	y int
}

func CreateRandomLocation() Location {
	x := rand.Intn(maxX) + minX
	y := rand.Intn(maxY) + minY
	location, err := NewLocation(x, y)
	if err != nil {
		panic(err)
	}
	return location
}

func NewLocation(x, y int) (Location, error) {
	if (x < minX || y < minY) || (x > maxX || y > maxY) {
		return Location{}, ErrInvalidLocation
	}

	return Location{
		x: x,
		y: y,
	}, nil
}

func (l Location) DistanceTo(other Location) int {
	return int(math.Abs(float64(l.x-other.x)) + math.Abs(float64(l.y-other.y)))
}

func (l Location) Equals(other Location) bool {
	return l == other
}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func MaxLocation() Location {
	return Location{
		x: maxX,
		y: maxY,
	}
}

func MinLocation() Location {
	return Location{
		x: minX,
		y: minY,
	}
}
