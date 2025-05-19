package ports

import (
	"context"

	"github.com/delivery/internal/core/domain/model/kernel"
)

type GeoServiceClient interface {
	GetLocation(ctx context.Context, street string) (kernel.Location, error)
}
