package shared

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents something significant that happened in the domain
// WHY: Enables event-driven architecture and decoupling between aggregates
type DomainEvent interface {
    EventID() string
    EventName() string
    AggregateID() string
    AggregateType() string
    OccurredAt() time.Time
}

// BaseEvent provides common event fields
// WHAT: Embedded in all domain events to avoid repetition
type BaseEvent struct {
    ID         string    `json:"event_id"`
    OccurredOn time.Time `json:"occurred_at"`
}

func NewBaseEvent() BaseEvent {
    return BaseEvent{
        ID:         uuid.New().String(),
        OccurredOn: time.Now(),
    }
}

func (e BaseEvent) EventID() string {
    return e.ID
}

func (e BaseEvent) OccurredAt() time.Time {
    return e.OccurredOn
}
