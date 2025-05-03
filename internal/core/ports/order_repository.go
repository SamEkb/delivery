package ports

import (
	"context"

	"github.com/delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Add(ctx context.Context, order *order.Order) error
	Update(ctx context.Context, order *order.Order) error
	Get(ctx context.Context, orderID uuid.UUID) (*order.Order, error)
	GetFirstInStatusCreate(ctx context.Context) (*order.Order, error)
	GetAllInStatusAssigned(ctx context.Context) ([]*order.Order, error)
}
