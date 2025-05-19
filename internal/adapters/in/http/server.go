package http

import (
	"github.com/delivery/internal/core/application/usecases/commands"
	"github.com/delivery/internal/core/application/usecases/queries"
	"github.com/delivery/internal/generated/servers"
	"github.com/delivery/internal/pkg/errs"
)

var _ servers.ServerInterface = (*Server)(nil)

type Server struct {
	assignOrder             commands.AssignOrderHandler
	createOrder             commands.CreateOrderHandler
	createCourier           commands.CreateCourierHandler
	getAllCouriers          queries.GetAllCouriersHandler
	getAllUncompletedOrders queries.GetAllUncompletedOrdersHandler
}

func NewServer(
	assignOrder commands.AssignOrderHandler,
	createOrder commands.CreateOrderHandler,
	createCourier commands.CreateCourierHandler,
	getAllCouriers queries.GetAllCouriersHandler,
	getAllUncompletedOrders queries.GetAllUncompletedOrdersHandler,
) (*Server, error) {
	if assignOrder == nil {
		return nil, errs.NewValueIsRequiredError("assign order handler")
	}
	if createOrder == nil {
		return nil, errs.NewValueIsRequiredError("create order handler")
	}
	if createCourier == nil {
		return nil, errs.NewValueIsRequiredError("create courier handler")
	}
	if getAllCouriers == nil {
		return nil, errs.NewValueIsRequiredError("get all couriers handler")
	}
	if getAllUncompletedOrders == nil {
		return nil, errs.NewValueIsRequiredError("get all uncompleted orders handler")
	}

	return &Server{
		assignOrder:             assignOrder,
		createOrder:             createOrder,
		createCourier:           createCourier,
		getAllCouriers:          getAllCouriers,
		getAllUncompletedOrders: getAllUncompletedOrders,
	}, nil
}
