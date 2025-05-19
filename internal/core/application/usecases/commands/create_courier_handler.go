package commands

import (
	"context"

	"github.com/delivery/internal/core/domain/model/courier"
	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/errs"
)

type CreateCourierHandler interface {
	Handle(context.Context, CreateCourierCommand) error
}

var _ CreateCourierHandler = &createCourierCommandHandler{}

type createCourierCommandHandler struct {
	unitOfWork ports.UnitOfWork
}

func NewCreateCourierHandler(
	unitOfWork ports.UnitOfWork) (CreateCourierHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &createCourierCommandHandler{
		unitOfWork: unitOfWork,
	}, nil
}

func (ch *createCourierCommandHandler) Handle(ctx context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errs.NewValueIsRequiredError("add address command")
	}

	location := kernel.CreateRandomLocation()
	courierAgg, err := courier.NewCourier(command.Name(), command.Speed(), location)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.CourierRepository().Add(ctx, courierAgg)
	if err != nil {
		return err
	}
	return nil
}
