package commands

import (
	"context"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/errs"
)

type CreateOrderHandler interface {
	Handle(ctx context.Context, command *CreateOrderCommand) error
}

type addCreateOrderHandler struct {
	uow ports.UnitOfWork
}

func NewAddCreateOrderHandler(uow ports.UnitOfWork) (CreateOrderHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("unit of work")
	}

	return &addCreateOrderHandler{
		uow: uow,
	}, nil
}

func (h *addCreateOrderHandler) Handle(ctx context.Context, command *CreateOrderCommand) error {
	if !command.IsValid() {
		return errs.NewValidationError("command", "create order command is invalid")
	}

	orderAgg, err := h.uow.OrderRepository().Get(ctx, command.OrderID())
	if err != nil && !errs.IsNotFound(err) {
		return errs.NewDatabaseError("get", "order", err)
	}
	if orderAgg != nil {
		return errs.NewConflictError("order", command.OrderID().String(), "order already exists")
	}

	location := kernel.CreateRandomLocation()
	newOrder, err := order.NewOrder(command.OrderID(), location, command.Volume())
	if err != nil {
		return errs.NewBusinessErrorWithCause("create order", "failed to create order domain object", err)
	}

	err = h.uow.OrderRepository().Add(ctx, newOrder)
	if err != nil {
		return errs.NewDatabaseError("add", "order", err)
	}

	return nil
}
