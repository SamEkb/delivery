package geo

import (
	"context"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/generated/clients/geosrv/geopb"
)

func (c *Client) GetLocation(ctx context.Context, street string) (kernel.Location, error) {
	req := &geopb.GetGeolocationRequest{Street: street}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	res, err := c.client.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}

	location, err := kernel.NewLocation(
		int(res.GetLocation().GetX()),
		int(res.GetLocation().GetY()),
	)
	if err != nil {
		return kernel.Location{}, err
	}

	return location, nil
}
