package queries

import (
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/errs"
)

type GetAllUncompletedOrdersHandler interface {
	Handle(query GetAllUncompletedOrdersQuery) (GetAllUncompletedOrdersResponse, error)
}

type getAllUncompletedOrdersHandler struct {
	uow ports.UnitOfWork
}

func NewGetAllUncompletedOrdersHandler(uow ports.UnitOfWork) (GetAllUncompletedOrdersHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("unit of work")
	}
	return &getAllUncompletedOrdersHandler{
		uow: uow,
	}, nil
}

func (h *getAllUncompletedOrdersHandler) Handle(query GetAllUncompletedOrdersQuery) (GetAllUncompletedOrdersResponse, error) {
	if !query.IsValid() {
		return GetAllUncompletedOrdersResponse{}, errs.NewValidationError("query", "get all uncompleted orders query is invalid")
	}

	var orders []OrderResponse
	result := h.uow.Db().Where("status IN ?", []string{order.Created.String(), order.Assigned.String()}).
		Find(&orders)

	if result.Error != nil {
		return GetAllUncompletedOrdersResponse{}, errs.NewDatabaseError("get", "orders", result.Error)
	}

	return GetAllUncompletedOrdersResponse{
		Orders: orders,
	}, nil
}
