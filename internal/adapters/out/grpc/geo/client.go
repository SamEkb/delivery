package geo

import (
	"time"

	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/generated/clients/geosrv/geopb"
	"github.com/delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ ports.GeoServiceClient = (*Client)(nil)

type Client struct {
	conn    *grpc.ClientConn
	client  geopb.GeoClient
	timeout time.Duration
}

func NewGeoClient(address string, timeout time.Duration) (*Client, error) {
	if address == "" {
		return nil, errs.NewValueIsRequiredError("address")
	}
	if timeout == 0 {
		return nil, errs.NewValueIsRequiredError("timeout")
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := geopb.NewGeoClient(conn)

	return &Client{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
