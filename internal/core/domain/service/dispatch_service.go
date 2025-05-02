package service

import (
	"errors"

	"github.com/delivery/internal/core/domain/model/courier"
	"github.com/delivery/internal/core/domain/model/order"
)

var ErrInvalidOrder = errors.New("order must not be nil and must be in Completed status")
var ErrInvalidCouriers = errors.New("couriers must not be nil")
var ErrCourierNotFound = errors.New("courier not found")

type DispatchService interface {
	Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

type dispatchService struct {
}

func NewDispatchService() DispatchService {
	return &dispatchService{}
}

func (d *dispatchService) Dispatch(orderParam *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if orderParam == nil || orderParam.Status() != order.Created {
		return nil, ErrInvalidOrder
	}

	if len(couriers) == 0 {
		return nil, ErrInvalidCouriers
	}

	var minTime float64
	var bestCourier *courier.Courier

	for _, c := range couriers {
		currentTime := c.CalculateTimeToLocation(orderParam.Location())

		canTake, err := c.CanTakeOrder(orderParam)
		if err != nil {
			return nil, err
		}

		if canTake && (currentTime < minTime || minTime == 0) {
			minTime = currentTime
			bestCourier = c
		}
	}

	if bestCourier == nil {
		return nil, ErrCourierNotFound
	}

	courierID := bestCourier.ID()

	if err := orderParam.Assign(&courierID); err != nil {
		return nil, err
	}

	if err := bestCourier.TakeOrder(orderParam); err != nil {
		return nil, err
	}

	return bestCourier, nil
}
