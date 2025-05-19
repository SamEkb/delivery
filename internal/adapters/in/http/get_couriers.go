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

func (s *Server) GetCouriers(ctx echo.Context) error {
	query, err := queries.NewGetAllCouriersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	result, err := s.getAllCouriers.Handle(*query)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
		return problems.NewBadRequest(err.Error())
	}

	var couriers []servers.Courier
	for _, courier := range result.Couriers {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		courier := servers.Courier{
			Id:       courier.ID,
			Name:     courier.Name,
			Location: location,
		}

		couriers = append(couriers, courier)
	}

	return ctx.JSON(http.StatusOK, couriers)
}
