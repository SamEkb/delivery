package commands

import (
	"context"

	"github.com/delivery/internal/core/domain/service"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/errs"
)

type AssignOrderHandler interface {
	Handle(ctx context.Context, command *AssignOrderCommand) error
}

type assignOrderHandler struct {
	uow        ports.UnitOfWork
	dispatcher service.DispatchService
}

func NewAssignOrderHandler(uow ports.UnitOfWork, dispatcher service.DispatchService) (AssignOrderHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("unit of work")
	}

	if dispatcher == nil {
		return nil, errs.NewValueIsRequiredError("dispatcher service")
	}

	return &assignOrderHandler{
		uow:        uow,
		dispatcher: dispatcher,
	}, nil
}

func (h *assignOrderHandler) Handle(ctx context.Context, command *AssignOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValidationError("command", "assign order command is invalid")
	}

	createdOrder, err := h.uow.OrderRepository().GetFirstInStatusCreate(ctx)
	if err != nil {
		if errs.IsNotFound(err) {
			return errs.NewNotFoundError("order", "in created status")
		}
		return errs.NewDatabaseError("get", "order", err)
	}

	couriers, err := h.uow.CourierRepository().GetAllAvailable(ctx)
	if err != nil {
		return errs.NewDatabaseError("get", "available couriers", err)
	}
	if len(couriers) == 0 {
		return errs.NewBusinessError("assign order", "no available couriers found")
	}

	courier, err := h.dispatcher.Dispatch(createdOrder, couriers)
	if err != nil {
		return errs.NewBusinessErrorWithCause("dispatch order", "failed to assign order to courier", err)
	}

	h.uow.Begin(ctx)

	if err := h.uow.OrderRepository().Update(ctx, createdOrder); err != nil {
		return errs.NewDatabaseError("update", "order", err)
	}
	if err := h.uow.CourierRepository().Update(ctx, courier); err != nil {
		return errs.NewDatabaseError("update", "courier", err)
	}

	if err = h.uow.Commit(ctx); err != nil {
		return errs.NewDatabaseError("commit", "transaction", err)
	}

	return nil
}
