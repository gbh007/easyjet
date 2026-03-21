package entity

type EventType int

const (
	EventCreated EventType = iota
	EventUpdated
	EventDeleted
	EventRequireAppRestart
)

type ProjectEvent struct {
	Type      EventType
	ProjectID uint
	Schedule  string
	Enabled   bool
}

type AppEvent struct {
	Type      EventType
	ProjectID uint
}
