package interfaces

import (
	"context"

	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
)

// EventPublisher defines how application publishes domain events
// WHY: Application layer needs to publish events without knowing infrastructure details
// WHERE: Injected into command handlers to publish events after operations
type EventPublisher interface {
    Publish(ctx context.Context, events ...shared.DomainEvent) error
}
