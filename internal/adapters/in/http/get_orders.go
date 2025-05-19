package http

import (
	"errors"
	"net/http"

	"github.com/delivery/internal/adapters/in/http/problems"
	"github.com/delivery/internal/core/application/usecases/queries"
	"github.com/delivery/internal/generated/servers"
	"github.com/delivery/internal/pkg/errs"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetOrders(ctx echo.Context) error {
	query, err := queries.NewGetAllUncompletedOrdersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	result, err := s.getAllUncompletedOrders.Handle(*query)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
		return problems.NewBadRequest(err.Error())
	}

	var orders []servers.Order
	for _, order := range result.Orders {
		location := servers.Location{
			X: order.Location.X,
			Y: order.Location.Y,
		}

		order := servers.Order{
			Id:       order.ID,
			Location: location,
		}

		orders = append(orders, order)
	}

	return ctx.JSON(http.StatusOK, orders)
}
