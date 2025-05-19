package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/delivery/internal/core/domain/model/kernel"
	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateOrderHandler_Handle(t *testing.T) {
	ctx := context.Background()

	mustCreateLocation := func(x, y int) kernel.Location {
		loc, err := kernel.NewLocation(x, y)
		if err != nil {
			panic(err)
		}
		return loc
	}

	type args struct {
		ctx     context.Context
		command *CreateOrderCommand
	}

	tests := map[string]struct {
		args    args
		wantErr bool
		deps    func(t *testing.T) (ports.UnitOfWork, ports.GeoServiceClient)
	}{
		"create success": {
			args: args{
				ctx: ctx,
				command: func() *CreateOrderCommand {
					orderID := uuid.New()
					cmd, _ := NewCreateOrderCommand(orderID, "street", 1)
					return cmd
				}(),
			},
			wantErr: false,
			deps: func(t *testing.T) (ports.UnitOfWork, ports.GeoServiceClient) {
				uow := mocks.NewUnitOfWork(t)
				orderRepo := mocks.NewOrderRepository(t)
				geoClient := mocks.NewGeoServiceClient(t)

				uow.EXPECT().OrderRepository().Return(orderRepo)

				orderRepo.EXPECT().
					Get(ctx, mock.MatchedBy(func(id uuid.UUID) bool {
						return id != uuid.Nil
					})).
					Return(nil, nil)

				geoClient.EXPECT().GetLocation(ctx, "street").Return(
					mustCreateLocation(1, 1),
					nil,
				)

				orderRepo.EXPECT().
					Add(ctx, mock.MatchedBy(func(o *order.Order) bool {
						return o != nil
					})).
					Return(nil)

				return uow, geoClient
			},
		},
		"order already exists": {
			args: args{
				ctx: ctx,
				command: func() *CreateOrderCommand {
					orderID := uuid.New()
					cmd, _ := NewCreateOrderCommand(orderID, "street", 1)
					return cmd
				}(),
			},
			wantErr: true,
			deps: func(t *testing.T) (ports.UnitOfWork, ports.GeoServiceClient) {
				uow := mocks.NewUnitOfWork(t)
				orderRepo := mocks.NewOrderRepository(t)
				geoClient := mocks.NewGeoServiceClient(t)

				existingOrderID := uuid.New()
				location := mustCreateLocation(1, 1)
				existingOrder := order.RestoreOrder(existingOrderID, nil, location, 1, order.Created)

				uow.EXPECT().OrderRepository().Return(orderRepo)

				orderRepo.EXPECT().
					Get(ctx, mock.MatchedBy(func(id uuid.UUID) bool {
						return id != uuid.Nil
					})).
					Return(existingOrder, nil)

				return uow, geoClient
			},
		},
		"error adding order": {
			args: args{
				ctx: ctx,
				command: func() *CreateOrderCommand {
					orderID := uuid.New()
					cmd, _ := NewCreateOrderCommand(orderID, "street", 1)
					return cmd
				}(),
			},
			wantErr: true,
			deps: func(t *testing.T) (ports.UnitOfWork, ports.GeoServiceClient) {
				uow := mocks.NewUnitOfWork(t)
				orderRepo := mocks.NewOrderRepository(t)
				geoClient := mocks.NewGeoServiceClient(t)

				uow.EXPECT().OrderRepository().Return(orderRepo)

				orderRepo.EXPECT().
					Get(ctx, mock.Anything).
					Return(nil, nil)

				geoClient.EXPECT().GetLocation(ctx, "street").Return(
					mustCreateLocation(1, 1),
					nil,
				)

				orderRepo.EXPECT().
					Add(ctx, mock.Anything).
					Return(errors.New("database error"))

				return uow, geoClient
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			uow, geo := tt.deps(t)
			handler, err := NewAddCreateOrderHandler(uow, geo)
			assert.NoError(t, err)

			err = handler.Handle(tt.args.ctx, tt.args.command)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
