package commands

import (
	"context"

	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/errs"
)

type MoveCourierHandler interface {
	Handle(ctx context.Context, command *MoveCourierCommand) error
}

type moveCourierHandler struct {
	uow ports.UnitOfWork
}

func NewMoveCourierHandler(uow ports.UnitOfWork) (MoveCourierHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("unit of work")
	}

	return &moveCourierHandler{
		uow: uow,
	}, nil
}

func (h *moveCourierHandler) Handle(ctx context.Context, command *MoveCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValidationError("command", "move courier command is invalid")
	}

	assignedOrders, err := h.uow.OrderRepository().GetAllInStatusAssigned(ctx)
	if err != nil {
		return errs.NewDatabaseError("get", "assigned orders", err)
	}

	h.uow.Begin(ctx)
	for _, order := range assignedOrders {
		courierID := order.CourierID()
		courier, err := h.uow.CourierRepository().Get(ctx, *courierID)
		if err != nil {
			return errs.NewDatabaseError("get", "courier", err)
		}

		if courier.Location().Equals(order.Location()) {
			if err := order.Complete(); err != nil {
				return errs.NewBusinessError("complete order", err.Error())
			}
			if err := courier.CompleteOrder(order); err != nil {
				return errs.NewBusinessError("complete order", err.Error())
			}
		} else {
			if err := courier.Move(order.Location()); err != nil {
				return errs.NewBusinessError("move courier", err.Error())
			}
		}

		if err := h.uow.CourierRepository().Update(ctx, courier); err != nil {
			return errs.NewDatabaseError("update", "courier", err)
		}
		if err := h.uow.OrderRepository().Update(ctx, order); err != nil {
			return errs.NewDatabaseError("update", "order", err)
		}
	}

	if err = h.uow.Commit(ctx); err != nil {
		return errs.NewDatabaseError("commit", "transaction", err)
	}

	return nil
}
