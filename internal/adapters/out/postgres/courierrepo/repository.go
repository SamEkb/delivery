package courierrepo

import (
	"context"

	"github.com/delivery/internal/core/domain/model/courier"
	"github.com/delivery/internal/core/ports"
	"github.com/delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CourierRepository = &Repository{}

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

func (r *Repository) Add(ctx context.Context, courier *courier.Courier) error {
	r.uow.Track(courier)

	dto := DomainToDto(courier)

	// check if we inside other tx
	isInTx := r.uow.InTx()
	if !isInTx {
		// if not, create own tx
		r.uow.Begin(ctx)
	}
	tx := r.uow.Tx()

	if err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(dto).Error; err != nil {
		return errs.NewDatabaseError("create", "courier", err)
	}

	// if inside other tx we not fix this one
	if !isInTx {
		if err := r.uow.Commit(ctx); err != nil {
			return errs.NewDatabaseError("commit", "transaction", err)
		}
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, courier *courier.Courier) error {
	r.uow.Track(courier)

	dto := DomainToDto(courier)

	// check if we inside other tx
	isInTx := r.uow.InTx()
	if !isInTx {
		// if not, create own tx
		r.uow.Begin(ctx)
	}
	tx := r.uow.Tx()

	if err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error; err != nil {
		return errs.NewDatabaseError("update", "courier", err)
	}

	// if inside other tx we not fix this one
	if !isInTx {
		if err := r.uow.Commit(ctx); err != nil {
			return errs.NewDatabaseError("commit", "transaction", err)
		}
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, courierID uuid.UUID) (*courier.Courier, error) {
	dto := CourierDto{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, courierID)

	if result.Error != nil {
		return nil, errs.NewDatabaseError("get", "courier", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errs.NewNotFoundError("courier", courierID.String())
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []CourierDto

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dtos)

	if result.Error != nil {
		return nil, errs.NewDatabaseError("get", "couriers", result.Error)
	}

	couriers := make([]*courier.Courier, len(dtos))
	for i, dto := range dtos {
		couriers[i] = DtoToDomain(dto)
	}

	return couriers, nil
}

func (r *Repository) GetAllAvailable(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []CourierDto

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Joins(`
        LEFT JOIN storage_places sp ON 
            couriers.id = sp.courier_id AND 
            sp.order_id IS NOT NULL
    `).
		Where("sp.id IS NULL").
		Preload(clause.Associations).
		Find(&dtos)

	if result.Error != nil {
		return nil, errs.NewDatabaseError("get", "available couriers", result.Error)
	}

	couriers := make([]*courier.Courier, len(dtos))
	for i, dto := range dtos {
		couriers[i] = DtoToDomain(dto)
	}

	return couriers, nil
}

func (r *Repository) getTxOrDb() *gorm.DB {
	if tx := r.uow.Tx(); tx != nil {
		return tx
	}
	return r.uow.Db()
}
