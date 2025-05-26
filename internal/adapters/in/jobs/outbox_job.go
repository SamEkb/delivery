package jobs

import (
	"context"
	"time"

	"github.com/delivery/internal/adapters/out/postgres/outboxrepo"
	"github.com/delivery/internal/pkg/ddd"
	"github.com/delivery/internal/pkg/errs"
	"github.com/delivery/internal/pkg/outbox"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &OutboxJob{}

type OutboxJob struct {
	outboxRepository outboxrepo.OutboxRepository
	registry         outbox.EventRegistry
	mediatr          ddd.Mediatr
}

func NewOutboxJob(
	outboxRepository outboxrepo.OutboxRepository,
	registry outbox.EventRegistry,
	mediatr ddd.Mediatr,
) (*OutboxJob, error) {
	if outboxRepository == nil {
		return nil, errs.NewValueIsRequiredError("outbox repository")
	}
	if registry == nil {
		return nil, errs.NewValueIsRequiredError("event registry")
	}
	if mediatr == nil {
		return nil, errs.NewValueIsRequiredError("mediatr")
	}

	return &OutboxJob{
		outboxRepository: outboxRepository,
		registry:         registry,
		mediatr:          mediatr,
	}, nil
}

func (outboxJob *OutboxJob) Run() {
	ctx := context.Background()

	events, err := outboxJob.outboxRepository.GetNotPublishedEvents()
	if err != nil {
		log.Error(err)
	}

	for _, event := range events {
		domainEvent, err := outboxJob.registry.DecodeDomainEvent(event)
		if err != nil {
			log.Error(err)
			continue
		}

		if err := outboxJob.mediatr.Publish(ctx, domainEvent); err != nil {
			log.Error(err)
			continue
		}

		now := time.Now().UTC()
		event.ProcessedAtUtc = &now
		if err := outboxJob.outboxRepository.Update(nil, event); err != nil {
			log.Error(err)
			continue
		}
	}
}
