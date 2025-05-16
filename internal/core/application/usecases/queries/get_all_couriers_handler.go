package queries

import (
	"errors"

	"gorm.io/gorm"
)

var ErrInvalidDBValue = errors.New("database must not be nil")
var ErrInvalidCommand = errors.New("command is invalid")

type GetAllCouriersHandler interface {
	Handle(query GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

type getAllCouriersHandler struct {
	db *gorm.DB
}

func NewGetAllCouriersHandler(db *gorm.DB) (GetAllCouriersHandler, error) {
	if db == nil {
		return nil, ErrInvalidDBValue
	}
	return &getAllCouriersHandler{
		db: db,
	}, nil
}

func (h *getAllCouriersHandler) Handle(query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if !query.IsValid() {
		return GetAllCouriersResponse{}, ErrInvalidCommand
	}
	var couriers []CourierResponse
	if err := h.db.Find(&couriers).Error; err != nil {
		return GetAllCouriersResponse{}, err
	}

	return GetAllCouriersResponse{
		Couriers: couriers,
	}, nil
}
