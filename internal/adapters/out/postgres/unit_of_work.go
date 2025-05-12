package postgres

import (
	"context"
	"errors"

	"github.com/delivery/internal/adapters/out/postgres/courierrepo"
	"github.com/delivery/internal/adapters/out/postgres/orderrepo"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/ddd"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var ErrDBMustBeNotNil = errors.New("database must be not nil")
var ErrInvalidTx = errors.New("cant commit transaction must be not nil")

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	tx                *gorm.DB
	db                *gorm.DB
	trackedAggregates []ddd.AggregateRoot
	courierRepository ports.CourierRepository
	orderRepository   ports.OrderRepository
}

func NewUnitOfWork(db *gorm.DB) (*UnitOfWork, error) {
	if db == nil {
		return nil, ErrDBMustBeNotNil
	}

	uow := &UnitOfWork{
		db: db,
	}

	courierRepo, err := courierrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.courierRepository = courierRepo

	orderRepo, err := orderrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.orderRepository = orderRepo

	return uow, nil
}

func (uow *UnitOfWork) Tx() *gorm.DB {
	return uow.tx
}

func (uow *UnitOfWork) Db() *gorm.DB {
	return uow.db
}

func (uow *UnitOfWork) InTx() bool {
	return uow.tx != nil
}

func (uow *UnitOfWork) Track(agg ddd.AggregateRoot) {
	uow.trackedAggregates = append(uow.trackedAggregates, agg)
}

func (uow *UnitOfWork) Begin(ctx context.Context) {
	uow.tx = uow.db.WithContext(ctx).Begin()
}

func (uow *UnitOfWork) Commit(ctx context.Context) error {
	if uow.tx == nil {
		return ErrInvalidTx
	}

	commited := false
	defer func() {
		if !commited {
			if err := uow.tx.Rollback().Error; err != nil {
				panic(err)
			}
			uow.clearTx()
		}
	}()

	if err := uow.tx.WithContext(ctx).Commit().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
		log.Error(err)
	}

	commited = true
	uow.clearTx()

	return nil
}

func (uow *UnitOfWork) CourierRepository() ports.CourierRepository {
	return uow.courierRepository
}

func (uow *UnitOfWork) OrderRepository() ports.OrderRepository {
	return uow.orderRepository
}

func (uow *UnitOfWork) clearTx() {
	uow.tx = nil
	uow.trackedAggregates = nil
}
