package entity

type EventType int

const (
	EventCreated EventType = iota
	EventUpdated
	EventDeleted
)

type ProjectEvent struct {
	Type      EventType
	ProjectID uint
	Schedule  string
	Enabled   bool
}
