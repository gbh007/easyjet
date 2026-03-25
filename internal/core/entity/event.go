package entity

type EventType int

const (
	EventProjectCreated EventType = iota
	EventProjectUpdated
	EventProjectDeleted
	EventRequireAppRestart
)

type Event struct {
	Type      EventType
	ProjectID uint
	RunID     uint
}
