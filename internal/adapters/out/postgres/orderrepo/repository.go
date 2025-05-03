package orderrepo

import (
	"context"
	"errors"

	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInvalidUOWValue = errors.New("unit of work must not be nil")
var ErrNoRecordsFound = errors.New("no records found")

var _ ports.OrderRepository = &Repository{}

type Repository struct {
	uow ports.UnitOfWork
}

func NewRepository(uow ports.UnitOfWork) (*Repository, error) {
	if uow == nil {
		return nil, ErrInvalidUOWValue
	}
	return &Repository{
		uow: uow,
	}, nil
}

func (r *Repository) Add(ctx context.Context, order *order.Order) error {
	r.uow.Track(order)

	dto := DomainToDto(order)

	// check if we inside other tx
	isInTx := r.uow.InTx()
	if !isInTx {
		// if not, create own tx
		r.uow.Begin(ctx)
	}
	tx := r.uow.Tx()

	if err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(dto).Error; err != nil {
		return err
	}

	// if inside other tx we not fix this one
	if !isInTx {
		if err := r.uow.Commit(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, order *order.Order) error {
	r.uow.Track(order)

	dto := DomainToDto(order)

	// check if we inside other tx
	isInTx := r.uow.InTx()
	if !isInTx {
		// if not, create own tx
		r.uow.Begin(ctx)
	}
	tx := r.uow.Tx()

	if err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error; err != nil {
		return err
	}

	// if inside other tx we not fix this one
	if !isInTx {
		if err := r.uow.Commit(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, orderID uuid.UUID) (*order.Order, error) {
	dto := OrderDTO{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, orderID)
	if result.RowsAffected == 0 {
		return nil, nil
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetFirstInStatusCreate(ctx context.Context) (*order.Order, error) {
	dto := OrderDTO{}
	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		First(&dto, "status = ?", order.Created)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecordsFound
		}
		return nil, result.Error
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetAllInStatusAssigned(ctx context.Context) ([]*order.Order, error) {
	var dtos []OrderDTO

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dtos, "status = ?", order.Assigned)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrNoRecordsFound
	}

	aggregates := make([]*order.Order, len(dtos))
	for i, dto := range dtos {
		aggregates[i] = DtoToDomain(dto)
	}

	return aggregates, nil
}

func (r *Repository) getTxOrDb() *gorm.DB {
	if tx := r.uow.Tx(); tx != nil {
		return tx
	}
	return r.uow.Db()
}
