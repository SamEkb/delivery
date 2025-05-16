package queries

import (
	"github.com/delivery/internal/core/domain/model/order"
	"gorm.io/gorm"
)

type GetAllUncompletedOrdersHandler interface {
	Handle(query GetAllUncompletedOrdersQuery) (GetAllUncompletedOrdersResponse, error)
}

type getAllUncompletedOrdersHandler struct {
	db *gorm.DB
}

func NewGetAllUncompletedOrdersHandler(db *gorm.DB) (GetAllUncompletedOrdersHandler, error) {
	if db == nil {
		return nil, ErrInvalidDBValue
	}
	return &getAllUncompletedOrdersHandler{
		db: db,
	}, nil
}

func (h *getAllUncompletedOrdersHandler) Handle(query GetAllUncompletedOrdersQuery) (GetAllUncompletedOrdersResponse, error) {
	if !query.IsValid() {
		return GetAllUncompletedOrdersResponse{}, ErrInvalidCommand
	}

	var orders []OrderResponse
	result := h.db.Where("status IN ?", []string{order.Created.String(), order.Assigned.String()}).
		Find(&orders)

	if result.Error != nil {
		return GetAllUncompletedOrdersResponse{}, result.Error
	}

	return GetAllUncompletedOrdersResponse{
		Orders: orders,
	}, nil
}
