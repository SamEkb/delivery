package queries

import (
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/errs"
)

type GetAllCouriersHandler interface {
	Handle(query GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

type getAllCouriersHandler struct {
	uow ports.UnitOfWork
}

func NewGetAllCouriersHandler(uow ports.UnitOfWork) (GetAllCouriersHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("unit of work")
	}
	return &getAllCouriersHandler{
		uow: uow,
	}, nil
}

func (h *getAllCouriersHandler) Handle(query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if !query.IsValid() {
		return GetAllCouriersResponse{}, errs.NewValidationError("query", "get all couriers query is invalid")
	}
	var couriers []CourierResponse
	if err := h.uow.Db().Find(&couriers).Error; err != nil {
		return GetAllCouriersResponse{}, errs.NewDatabaseError("get", "couriers", err)
	}

	return GetAllCouriersResponse{
		Couriers: couriers,
	}, nil
}
