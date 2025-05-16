package commands

import (
	"context"

	"github.com/delivery/internal/core/ports"
)

type MoveCourierHandler interface {
	Handler(ctx context.Context, command *MoveCourierCommand) error
}

type moveCourierHandler struct {
	uow ports.UnitOfWork
}

func NewMoveCourierHandler(uow ports.UnitOfWork) (MoveCourierHandler, error) {
	if uow == nil {
		return nil, ErrInvalidUOWValue
	}

	return &moveCourierHandler{
		uow: uow,
	}, nil
}

func (h *moveCourierHandler) Handler(ctx context.Context, command *MoveCourierCommand) error {
	if !command.isValid {
		return ErrInvalidCommand
	}

	assignedOrders, err := h.uow.OrderRepository().GetAllInStatusAssigned(ctx)
	if err != nil {
		return err
	}

	h.uow.Begin(ctx)
	for _, order := range assignedOrders {
		courierID := order.CourierID()
		courier, err := h.uow.CourierRepository().Get(ctx, *courierID)
		if err != nil {
			return err
		}

		if courier.Location().Equals(order.Location()) {
			if err := order.Complete(); err != nil {
				return err
			}
			if err := courier.CompleteOrder(order); err != nil {
				return err
			}
		} else {
			if err := courier.Move(order.Location()); err != nil {
				return err
			}
		}

		if err := h.uow.CourierRepository().Update(ctx, courier); err != nil {
			return err
		}
		if err := h.uow.OrderRepository().Update(ctx, order); err != nil {
			return err
		}
	}

	if err = h.uow.Commit(ctx); err != nil {
		return err
	}

	return nil
}
