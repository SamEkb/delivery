package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID             uuid.UUID
	Name           string
	Payload        []byte
	OccurredAtUtc  time.Time
	ProcessedAtUtc *time.Time
}

func (Event) TableName() string {
	return "outbox"
}
