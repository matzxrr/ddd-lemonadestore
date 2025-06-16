package shared

// AggregateRoot is the base for all aggregate roots in the domain
// WHY: Provides common functionality for event sourcing and aggregate identification
// WHAT: Stores domain events that occurred during business operations
type AggregateRoot struct {
    events []DomainEvent
}

// Raise adds a domain event to the aggregate
// WHERE: Called within domain methods when significant business events occur
func (a *AggregateRoot) Raise(event DomainEvent) {
    a.events = append(a.events, event)
}

// PullEvents returns all raised events and clears the internal list
// WHY: Allows infrastructure layer to retrieve and publish events after persistence
func (a *AggregateRoot) PullEvents() []DomainEvent {
    events := a.events
    a.events = []DomainEvent{} // Clear after pulling
    return events
}
