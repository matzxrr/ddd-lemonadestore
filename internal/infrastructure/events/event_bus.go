package events

import (
	"context"
	"log"
	"sync"

	"github.com/matzxrr/ddd-lemonadestore/internal/application/interfaces"
	"github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
)

// EventHandler is a function that handles domain events
type EventHandler func(ctx context.Context, event shared.DomainEvent) error

// InMemoryEventBus is an in-memory implementation of event publishing
// WHY: Decouples event producers from consumers
// WHERE: Used by application layer to publish domain events
type InMemoryEventBus struct {
    mu       sync.RWMutex
    handlers map[string][]EventHandler
}

// NewInMemoryEventBus creates a new event bus
func NewInMemoryEventBus() *InMemoryEventBus {
    return &InMemoryEventBus{
        handlers: make(map[string][]EventHandler),
    }
}

// Subscribe registers a handler for an event type
// WHAT: Allows registration of multiple handlers per event
func (bus *InMemoryEventBus) Subscribe(eventName string, handler EventHandler) {
    bus.mu.Lock()
    defer bus.mu.Unlock()
    
    bus.handlers[eventName] = append(bus.handlers[eventName], handler)
}

// Publish sends events to all registered handlers
// WHY: Implements eventual consistency pattern
func (bus *InMemoryEventBus) Publish(ctx context.Context, events ...shared.DomainEvent) error {
    for _, event := range events {
        bus.mu.RLock()
        handlers := bus.handlers[event.EventName()]
        bus.mu.RUnlock()
        
        // Execute handlers asynchronously
        for _, handler := range handlers {
            go func(h EventHandler, e shared.DomainEvent) {
                if err := h(ctx, e); err != nil {
                    log.Printf("Event handler error for %s: %v", e.EventName(), err)
                }
            }(handler, event)
        }
    }
    
    return nil
}

// Ensure it implements the interface
var _ interfaces.EventPublisher = (*InMemoryEventBus)(nil)
