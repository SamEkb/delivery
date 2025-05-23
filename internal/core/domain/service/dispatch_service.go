package service

import (
	"github.com/delivery/internal/core/domain/model/courier"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/pkg/errs"
)

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
		return nil, errs.NewValidationError("order", "order is not created")
	}

	if couriers == nil || len(couriers) == 0 {
		return nil, errs.NewValidationError("couriers", "couriers not found")
	}

	bestCourier, err := findNearestSuitableCourier(orderParam, couriers)
	if err != nil {
		return nil, err
	}

	if bestCourier == nil {
		return nil, errs.NewValidationError("couriers", "no suitable couriers found")
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

func findNearestSuitableCourier(orderParam *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	var minTime float64
	var suitableCourier *courier.Courier
	for _, c := range couriers {
		currentTime := c.CalculateTimeToLocation(orderParam.Location())

		canTake, err := c.CanTakeOrder(orderParam)
		if err != nil {
			return nil, err
		}

		if canTake && (currentTime < minTime || minTime == 0) {
			minTime = currentTime
			suitableCourier = c
		}
	}
	return suitableCourier, nil
}
