package outboxrepo

import (
	"context"

	"github.com/delivery/internal/pkg/errs"
	"github.com/delivery/internal/pkg/outbox"
	"gorm.io/gorm"
)

type OutboxRepository interface {
	Update(ctx context.Context, message *outbox.Event) error
	GetNotPublishedEvents() ([]*outbox.Event, error)
}

var _ OutboxRepository = &repository{}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*repository, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &repository{
		db: db,
	}, nil
}

func (r *repository) Update(ctx context.Context, message *outbox.Event) error {
	if err := r.db.WithContext(ctx).Save(&message).Error; err != nil {
		return errs.NewDatabaseError("update", "outbox event", err)
	}

	return nil
}

func (r *repository) GetNotPublishedEvents() ([]*outbox.Event, error) {
	var events []*outbox.Event
	if err := r.db.WithContext(context.Background()).
		Order("occurred_at_utc ASC").
		Limit(20).
		Where("processed_at_utc IS NULL").
		Find(&events).Error; err != nil {
		return nil, errs.NewDatabaseError("get", "not published events", err)
	}
	return events, nil
}
