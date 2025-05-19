package http

import (
	"errors"
	"net/http"

	"github.com/delivery/internal/adapters/in/http/problems"
	"github.com/delivery/internal/core/application/usecases/commands"
	"github.com/delivery/internal/generated/servers"
	"github.com/delivery/internal/pkg/errs"
	"github.com/labstack/echo/v4"
)

func (s *Server) CreateCourier(ctx echo.Context) error {
	var courier servers.NewCourier
	if err := ctx.Bind(&courier); err != nil {
		return problems.NewBadRequest(err.Error())
	}

	command, err := commands.NewCreateCourierCommand(courier.Name, courier.Speed)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	if err = s.createCourier.Handle(ctx.Request().Context(), command); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return problems.NewNotFound(err.Error())
		}
		return problems.NewConflict(err.Error(), "/")

	}

	return ctx.JSON(http.StatusOK, nil)
}
