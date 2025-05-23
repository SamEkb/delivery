package eventhandlers

import (
	"context"

	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/ddd"
	"github.com/delivery/internal/pkg/errs"
)

type orderStatusChangedEventHandler struct {
	orderProducer ports.OrderProducer
}

func NewOrderStatusChangedEventHandler(orderProducer ports.OrderProducer) (ddd.EventHandler, error) {
	if orderProducer == nil {
		return nil, errs.NewValueIsRequiredError("order producer")
	}
	return &orderStatusChangedEventHandler{
		orderProducer: orderProducer,
	}, nil
}

func (h *orderStatusChangedEventHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	if err := h.orderProducer.Publish(ctx, event); err != nil {
		return err
	}

	return nil
}
