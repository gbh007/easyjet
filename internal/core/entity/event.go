package entity

import "time"

type EventType int

const (
	EventProjectCreated EventType = iota
	EventProjectUpdated
	EventProjectDeleted
	EventRunFinished
	EventRunGitFinished
	EventRunStageFinished
	EventRunRotateFinished
	EventRequireAppRestart
)

type Event struct {
	Type      EventType
	ProjectID uint
	RunID     uint
	Stage     int
	Err       error
	Duration  time.Duration
}
