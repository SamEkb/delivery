package http

import (
	"net/http"

	"github.com/delivery/internal/adapters/in/http/problems"
	"github.com/delivery/internal/core/application/usecases/commands"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) CreateOrder(ctx echo.Context) error {
	command, err := commands.NewCreateOrderCommand(uuid.New(), "Несуществующая", 1)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.createOrder.Handle(ctx.Request().Context(), command)
	if err != nil {
		return problems.NewConflict(err.Error(), "/")
	}

	return ctx.JSON(http.StatusOK, nil)
}
