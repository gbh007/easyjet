package scheduler

// EventType represents the type of scheduler event.
type EventType int

const (
	// EventCreated indicates a project was created with cron scheduling.
	EventCreated EventType = iota
	// EventUpdated indicates a project's cron settings were updated.
	EventUpdated
	// EventDeleted indicates a project was deleted.
	EventDeleted
)

// SchedulerEvent represents an event in the scheduler event queue.
type SchedulerEvent struct {
	// Type is the type of event (Created, Updated, Deleted).
	Type EventType
	// ProjectID is the ID of the project associated with this event.
	ProjectID uint
	// Schedule is the cron schedule expression (empty if deleted or no schedule).
	Schedule string
	// Enabled indicates whether cron scheduling is enabled for this project.
	Enabled bool
}

// EventQueue is a channel-based queue for scheduler events.
type EventQueue struct {
	events chan SchedulerEvent
}

// NewEventQueue creates a new event queue with the specified buffer size.
// FIXME(ai-shit): сделать более унифицировано.
func NewEventQueue(bufferSize int) *EventQueue {
	return &EventQueue{
		events: make(chan SchedulerEvent, bufferSize),
	}
}

// Publish sends an event to the queue.
// Returns an error if the queue is closed or full.
func (eq *EventQueue) Publish(event SchedulerEvent) error {
	eq.events <- event
	return nil
}

// Subscribe returns the read-only channel for consuming events.
func (eq *EventQueue) Subscribe() <-chan SchedulerEvent {
	return eq.events
}

// Close closes the event queue channel.
// Should be called when the application shuts down.
func (eq *EventQueue) Close() {
	close(eq.events)
}

// EventPublisher defines the interface for publishing scheduler events.
type EventPublisher interface {
	Publish(event SchedulerEvent) error
}
