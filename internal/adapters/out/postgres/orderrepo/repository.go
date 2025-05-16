package orderrepo

import (
	"context"

	"github.com/delivery/internal/core/domain/model/order"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.OrderRepository = &Repository{}

type Repository struct {
	uow ports.UnitOfWork
}

func NewRepository(uow ports.UnitOfWork) (*Repository, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("unit of work")
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
		return errs.NewDatabaseError("create", "order", err)
	}

	// if inside other tx we not fix this one
	if !isInTx {
		if err := r.uow.Commit(ctx); err != nil {
			return errs.NewDatabaseError("commit", "transaction", err)
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
		return errs.NewDatabaseError("update", "order", err)
	}

	// if inside other tx we not fix this one
	if !isInTx {
		if err := r.uow.Commit(ctx); err != nil {
			return errs.NewDatabaseError("commit", "transaction", err)
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

	if result.Error != nil {
		return nil, errs.NewDatabaseError("get", "order", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errs.NewNotFoundError("order", orderID.String())
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
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFoundError("order", "in created status")
		}
		return nil, errs.NewDatabaseError("get", "order", result.Error)
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetAllInStatusCreate(ctx context.Context) ([]*order.Order, error) {
	var dtos []OrderDTO

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dtos, "status = ?", order.Created)

	if result.Error != nil {
		return nil, errs.NewDatabaseError("get", "orders", result.Error)
	}

	aggregates := make([]*order.Order, len(dtos))
	for i, dto := range dtos {
		aggregates[i] = DtoToDomain(dto)
	}

	return aggregates, nil
}

func (r *Repository) GetAllInStatusAssigned(ctx context.Context) ([]*order.Order, error) {
	var dtos []OrderDTO

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dtos, "status = ?", order.Assigned)

	if result.Error != nil {
		return nil, errs.NewDatabaseError("get", "orders", result.Error)
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
