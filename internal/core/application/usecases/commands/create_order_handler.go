package commands

import (
	"context"
	"errors"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/core/ports"
)

var ErrInvalidUOWValue = errors.New("unit of work must not be nil")
var ErrInvalidCommand = errors.New("command is invalid")

type CreateOrderHandler interface {
	Handle(ctx context.Context, command *CreateOrderCommand) error
}

type addCreateOrderHandler struct {
	uow ports.UnitOfWork
}

func NewAddCreateOrderHandler(uow ports.UnitOfWork) (CreateOrderHandler, error) {
	if uow == nil {
		return nil, ErrInvalidUOWValue
	}

	return &addCreateOrderHandler{
		uow: uow,
	}, nil
}

func (h *addCreateOrderHandler) Handle(ctx context.Context, command *CreateOrderCommand) error {
	if !command.IsValid() {
		return ErrInvalidCommand
	}

	orderAgg, err := h.uow.OrderRepository().Get(ctx, command.OrderID())
	if err != nil {
		return err
	}
	if orderAgg != nil {
		return nil
	}

	location := kernel.CreateRandomLocation()
	newOrder, err := order.NewOrder(command.OrderID(), location, command.Volume())
	if err != nil {
		return err
	}

	err = h.uow.OrderRepository().Add(ctx, newOrder)
	if err != nil {
		return err
	}

	return nil
}
