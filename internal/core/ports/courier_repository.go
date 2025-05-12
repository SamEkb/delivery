package ports

import (
	"context"

	"github.com/delivery/internal/core/domain/model/courier"
	"github.com/google/uuid"
)

type CourierRepository interface {
	Add(ctx context.Context, courier *courier.Courier) error
	Update(ctx context.Context, courier *courier.Courier) error
	Get(ctx context.Context, courierID uuid.UUID) (*courier.Courier, error)
	GetAllAvailable(ctx context.Context) ([]*courier.Courier, error)
}
