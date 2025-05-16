package commands

import (
	"context"
	"errors"

	"github.com/delivery/internal/core/domain/service"
	"github.com/delivery/internal/core/ports"
)

var ErrNoAvailableCouriers = errors.New("no available couriers")
var ErrInvalidDispatcherValue = errors.New("dispatcher must not be nil")

type AssignOrderHandler interface {
	Handle(ctx context.Context, command *AssignOrderCommand) error
}

type assignOrderHandler struct {
	uow        ports.UnitOfWork
	dispatcher service.DispatchService
}

func NewAssignOrderHandler(uow ports.UnitOfWork, dispatcher service.DispatchService) (AssignOrderHandler, error) {
	if uow == nil {
		return nil, ErrInvalidUOWValue
	}

	if dispatcher == nil {
		return nil, ErrInvalidDispatcherValue
	}

	return &assignOrderHandler{
		uow:        uow,
		dispatcher: dispatcher,
	}, nil
}

func (h *assignOrderHandler) Handle(ctx context.Context, command *AssignOrderCommand) error {
	if !command.IsValid() {
		return ErrInvalidCommand
	}

	createdOrder, err := h.uow.OrderRepository().GetFirstInStatusCreate(ctx)
	if err != nil {
		return err
	}

	couriers, err := h.uow.CourierRepository().GetAllAvailable(ctx)
	if err != nil {
		return err
	}
	if len(couriers) == 0 {
		return ErrNoAvailableCouriers
	}

	courier, err := h.dispatcher.Dispatch(createdOrder, couriers)
	if err != nil {
		return err
	}

	h.uow.Begin(ctx)

	if err := h.uow.OrderRepository().Update(ctx, createdOrder); err != nil {
		return err
	}
	if err := h.uow.CourierRepository().Update(ctx, courier); err != nil {
		return err
	}

	if err = h.uow.Commit(ctx); err != nil {
		return err
	}

	return nil
}
